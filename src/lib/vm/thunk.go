package vm

import (
	"sync"
	"sync/atomic"
)

type thunkState int32

const (
	illegal thunkState = iota
	normal
	locked
	app
)

type Thunk struct {
	Result   Object
	function *Thunk
	// `args` is represented as a slice but not a List to let users optimize
	// their functions. If you want to define functions with arguments fully
	// lazy, just create a function which takes only a thunk of a List as a
	// argument.
	args      []*Thunk
	state     thunkState
	blackHole sync.WaitGroup
}

func Normal(o Object) *Thunk {
	if f, ok := o.(func(...*Thunk) Object); ok {
		o = NewLazyFunction(f)
	} else if f, ok := o.(func(...Object) Object); ok {
		o = NewStrictFunction(f)
	}

	return &Thunk{Result: o, state: normal}
}

func App(f *Thunk, args ...*Thunk) *Thunk {
	t := &Thunk{function: f, args: args, state: app}
	t.blackHole.Add(1)
	return t
}

func (t *Thunk) Eval() Object { // into WHNF
	if t.compareAndSwapState(app, locked) {
		for {
			// This algorithm to eliminate tail calls is too hacky.
			o := t.function.Eval()
			f, ok := o.(Callable)

			if !ok {
				t.Result = NotCallableError(o)
				break
			}

			t.Result = f.Call(t.args...)
			child, ok := t.Result.(*Thunk)

			if !ok {
				break
			}

			// child.function and child.args can be extracted safely only if the
			// child thunk has only 1 reference. (e.g. tail calls)
			t.function = child.function
			t.args = child.args
		}

		t.function = nil
		t.args = nil

		t.storeState(normal)

		t.blackHole.Done()
	} else {
		t.blackHole.Wait()
	}

	return t.Result
}

func (t *Thunk) compareAndSwapState(old, new thunkState) bool {
	return atomic.CompareAndSwapInt32((*int32)(&t.state), int32(old), int32(new))
}

func (t *Thunk) storeState(new thunkState) {
	atomic.StoreInt32((*int32)(&t.state), int32(new))
}
