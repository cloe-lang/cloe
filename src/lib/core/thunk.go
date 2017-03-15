package core

import (
	"sync"
	"sync/atomic"

	"github.com/raviqqe/tisp/src/lib/debug"
	"github.com/raviqqe/tisp/src/lib/util"
)

type thunkState int32

const (
	normal thunkState = iota
	app
)

// Thunk you all!
type Thunk struct {
	result    Value
	function  *Thunk
	args      Arguments
	state     thunkState
	blackHole sync.WaitGroup
	info      debug.Info
}

// Normal creates a thunk of a WHNF value as its result.
func Normal(v Value) *Thunk {
	checkValue("Normal's argument", v)
	return &Thunk{result: v, state: normal}
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

// Eval evaluates a thunk and returns a WHNF value.
func (t *Thunk) Eval() Value {
	if t.lock() {
		children := make([]*Thunk, 0, 1) // TODO: Best cap?

		for {
			v := t.moveFunction().Eval()

			if t.chainError(v) {
				break
			}

			f, ok := v.(callable)

			if !ok {
				t.result = NotCallableError(v).Eval()
				break
			}

			t.result = f.call(t.moveArguments())

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

		checkValue("Thunk.result", t.result)

		for _, child := range children {
			// TODO: Use children's debug informations, child.info?
			child.result = t.result
			child.finalize()
		}

		t.finalize()
	} else {
		t.blackHole.Wait()
	}

	checkValue("Thunk.result", t.result)

	return t.result
}

func (t *Thunk) lock() bool {
	return t.compareAndSwapState(app, normal)
}

func (t *Thunk) delegateEval() (*Thunk, Arguments, bool) {
	if t.lock() {
		return t.moveFunction(), t.moveArguments(), true
	}

	return nil, Arguments{}, false
}

func (t *Thunk) moveFunction() *Thunk {
	f := t.function
	t.function = nil
	return f
}

func (t *Thunk) moveArguments() Arguments {
	args := t.args
	t.args = Arguments{}
	return args
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

func (t *Thunk) chainError(v Value) bool {
	if e, ok := v.(ErrorType); ok {
		e.callTrace = append(e.callTrace, t.info)
		t.result = e
		return true
	}

	return false
}

func checkValue(s string, v Value) {
	if _, ok := v.(*Thunk); ok {
		util.Fail(s + " is *Thunk.")
	}
}
