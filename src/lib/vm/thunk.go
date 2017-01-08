package vm

import (
	"sync"
	"sync/atomic"
)

type thunkState uint32

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
	if f, ok := o.(func(...*Thunk) *Thunk); ok {
		o = Function(f)
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
		o := t.function.Eval()
		f, ok := o.(Callable)

		if ok {
			t.Result = f.Call(t.args...).Eval()
		} else {
			t.Result = NotCallableError(o).Eval()
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
	return atomic.CompareAndSwapUint32(
		(*uint32)(&t.state),
		uint32(old),
		uint32(new))
}

func (t *Thunk) storeState(new thunkState) {
	atomic.StoreUint32((*uint32)(&t.state), uint32(new))
}
