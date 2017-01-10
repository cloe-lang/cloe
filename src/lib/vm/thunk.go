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
	result   Object
	function *Thunk
	// `args` is represented as a slice but not a List to let users optimize
	// their functions. If you want to define functions with arguments fully
	// lazy, just create a function which takes only a thunk of a List as a
	// argument.
	args         []*Thunk
	state        thunkState
	blackHole    sync.WaitGroup
	trampoliners []*Thunk
}

func Normal(o Object) *Thunk {
	if f, ok := o.(func(...*Thunk) Object); ok {
		o = NewLazyFunction(f)
	} else if f, ok := o.(func(...Object) Object); ok {
		o = NewStrictFunction(f)
	}

	return &Thunk{result: o, state: normal}
}

func App(f *Thunk, args ...*Thunk) *Thunk {
	t := &Thunk{function: f, args: args, state: app, trampoliners: make([]*Thunk, 0)}
	t.blackHole.Add(1)
	return t
}

func (t *Thunk) EvalStrictly() Object { // return WHNF
	if t.compareAndSwapState(app, locked) {
		o := t.function.EvalStrictly()
		f, ok := o.(Callable)

		if ok {
			t.result = f.Call(t.args...)

			for {
				child, ok := t.result.(*Thunk)

				if !ok {
					break
				}

				t.result = child.Eval()
			}
		} else {
			t.result = NotCallableError(o)
		}

		t.finalize()
	} else {
		t.blackHole.Wait()
	}

	return t.result
}

func (t *Thunk) Eval() Object { // return WHNF or *Thunk
	if t.compareAndSwapState(app, locked) {
		o := t.function.EvalStrictly()
		f, ok := o.(Callable)

		if ok {
			t.result = f.Call(t.args...)
			child, ok := t.result.(*Thunk)

			if ok {
				child.addTrampoliner(t)
				return child
			}
		} else {
			t.result = NotCallableError(o)
		}

		t.finalize()
	} else {
		t.blackHole.Wait()
	}

	return t.result
}

func (t *Thunk) storeResult(o Object) {
	t.result = o
	t.finalize()
}

func (t *Thunk) finalize() {
	for _, tr := range t.trampoliners {
		tr.storeResult(t.result)
	}
	t.trampoliners = nil

	t.function = nil
	t.args = nil
	t.storeState(normal)
	t.blackHole.Done()
}

func (t *Thunk) addTrampoliner(tr *Thunk) {
	t.trampoliners = append(t.trampoliners, tr)
}

func (t *Thunk) compareAndSwapState(old, new thunkState) bool {
	return atomic.CompareAndSwapInt32((*int32)(&t.state), int32(old), int32(new))
}

func (t *Thunk) storeState(new thunkState) {
	atomic.StoreInt32((*int32)(&t.state), int32(new))
}
