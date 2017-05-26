package core

import (
	"sync"
	"sync/atomic"

	"github.com/tisp-lang/tisp/src/lib/debug"
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
	assertValueIsWHNF("Normal's argument", v)
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

// EvalAny evaluates a thunk and returns a pure or impure (output) value.
func (t *Thunk) EvalAny(isPure bool) Value {
	if t.lock() {
		children := make([]*Thunk, 0)

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
				t.result = child.EvalAny(isPure)
				t.chainError(t.result)
				break
			}

			children = append(children, child)
		}

		assertValueIsWHNF("Thunk.result", t.result)

		if _, ok := t.result.(OutputType); isPure && ok {
			t.result = ImpureFunctionError(t.result).Eval()
		} else if !isPure && !ok {
			t.result = NotOutputError(t.result).Eval()
		}

		for _, child := range children {
			// TODO: Use children's debug informations, child.info?
			child.result = t.result
			child.finalize()
		}

		t.finalize()
	} else {
		t.blackHole.Wait()
	}

	assertValueIsWHNF("Thunk.result", t.result)

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

func assertValueIsWHNF(s string, v Value) {
	if _, ok := v.(*Thunk); ok {
		panic(s + " is *Thunk")
	}
}

// Eval evaluates a pure value.
func (t *Thunk) Eval() Value {
	return t.EvalAny(true)
}

// EvalOutput evaluates an output expression.
func (t *Thunk) EvalOutput() Value {
	v := t.EvalAny(false)

	if err, ok := v.(ErrorType); ok {
		return err
	}

	return v.(OutputType).value.Eval()
}
