package core

import (
	"github.com/raviqqe/tisp/src/lib/debug"
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
	result    Object
	function  *Thunk
	args      Arguments
	state     thunkState
	blackHole sync.WaitGroup
	info      debug.Info
}

// Normal creates a thunk of a WHNF object as its result.
func Normal(o Object) *Thunk {
	checkObject("Normal's argument", o)
	return &Thunk{result: o, state: normal}
}

// App creates a thunk applying a function to arguments.
func App(f *Thunk, args Arguments) *Thunk {
	return AppWithInfo(f, args, debug.NewGoInfo(1))
}

// AppWithInfo is the same as App except that it stores debug information
// in the thunk.
func AppWithInfo(f *Thunk, args Arguments, i debug.Info) *Thunk {
	t := &Thunk{
		function: f,
		args:     args,
		state:    app,
		info:     i,
	}
	t.blackHole.Add(1)
	return t
}

// PApp is not PPap.
func PApp(f *Thunk, ps ...*Thunk) *Thunk {
	return AppWithInfo(f, NewPositionalArguments(ps...), debug.NewGoInfo(1))
}

// Eval evaluates a thunk and returns a WHNF object.
func (t *Thunk) Eval() Object {
	if t.lock() {
		children := make([]*Thunk, 0, 1) // TODO: best capacity?

		for {
			o := t.function.Eval()

			if t.chainError(o) {
				break
			}

			f, ok := o.(callable)

			if !ok {
				t.result = NotCallableError(o).Eval()
				break
			}

			t.result = f.call(t.args)

			if t.chainError(t.result) {
				break
			}

			child, ok := t.result.(*Thunk)

			if !ok {
				break
			}

			t.function, t.args, ok = child.delegateEval()

			if !ok {
				t.result = child.Eval()
				t.chainError(t.result)
				break
			}

			children = append(children, child)
		}

		checkObject("Thunk.result", t.result)

		for _, child := range children {
			// TODO: Use children's debug informations, child.info?
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

func (t *Thunk) chainError(o Object) bool {
	if e, ok := o.(ErrorType); ok {
		e.callTrace = append(e.callTrace, t.info)
		t.result = e
		return true
	}

	return false
}

func checkObject(s string, o Object) {
	if _, ok := o.(*Thunk); ok {
		panic(s + " is *Thunk.")
	}
}
