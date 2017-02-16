package vm

import (
	"sync"
	"sync/atomic"
)

type thunkState int32

const (
	// States of illegal and locked are unnecessary and provided only for
	// debuggability.
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
	args      Arguments
	state     thunkState
	blackHole sync.WaitGroup
}

func Normal(o Object) *Thunk {
	checkObject("Normal's argument", o)

	return &Thunk{result: o, state: normal}
}

func App(f *Thunk, args Arguments) *Thunk {
	t := &Thunk{function: f, args: args, state: app}
	t.blackHole.Add(1)
	return t
}

func PApp(f *Thunk, ps ...*Thunk) *Thunk {
	return App(f, NewPositionalArguments(ps...))
}

func (t *Thunk) Eval() Object { // return WHNF
	if t.lock() {
		children := make([]*Thunk, 0, 1) // TODO: best capacity?

		for {
			o := t.function.Eval()
			f, ok := o.(callable)

			if !ok {
				t.result = NotCallableError(o).Eval()
				break
			}

			t.result = f.call(t.args)
			child, ok := t.result.(*Thunk)

			if !ok {
				break
			}

			t.function, t.args, ok = child.delegateEval()

			if !ok {
				t.result = child.Eval()
				break
			}

			children = append(children, child)
		}

		checkObject("Thunk.result", t.result)

		for _, child := range children {
			child.result = t.result
			child.finalize()
		}

		t.finalize()
	} else {
		t.blackHole.Wait()
	}

	checkObject("Thunk.result", t.result)

	return t.result
}

func (t *Thunk) lock() bool {
	return t.compareAndSwapState(app, locked)
}

func (t *Thunk) delegateEval() (*Thunk, Arguments, bool) {
	if t.lock() {
		return t.function, t.args, true
	}

	return nil, Arguments{}, false
}

func (t *Thunk) finalize() {
	t.function = nil
	t.args = Arguments{}
	t.storeState(normal)
	t.blackHole.Done()
}

func (t *Thunk) compareAndSwapState(old, new thunkState) bool {
	return atomic.CompareAndSwapInt32((*int32)(&t.state), int32(old), int32(new))
}

func (t *Thunk) storeState(new thunkState) {
	atomic.StoreInt32((*int32)(&t.state), int32(new))
}

func checkObject(s string, o Object) {
	if _, ok := o.(*Thunk); ok {
		panic(s + " is *Thunk.")
	}
}
