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
	spinLock
)

var identity = NewLazyFunction(
	NewSignature([]string{"x"}, nil, "", nil, nil, ""),
	func(ts ...*Thunk) Value { return ts[0] })

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
func (t *Thunk) EvalAny(pure bool) Value {
	if t.lock(normal) {
		for {
			v := t.swapFunction(nil).Eval()

			if t.chainError(v) {
				break
			}

			f, ok := v.(callable)

			if !ok {
				t.result = NotCallableError(v).Eval()
				break
			}

			t.result = f.call(t.swapArguments(Arguments{}))

			if t.chainError(t.result) {
				break
			}

			child, ok := t.result.(*Thunk)

			if !ok {
				break
			}

			t.function, t.args, ok = child.delegateEval(t)

			if !ok {
				t.result = child.EvalAny(pure)
				t.chainError(t.result)
				break
			}
		}

		assertValueIsWHNF("Thunk.result", t.result)

		if _, impure := t.result.(OutputType); pure && impure {
			t.result = ImpureFunctionError(t.result).Eval()
		} else if !pure && !impure {
			t.result = NotOutputError(t.result).Eval()
		}

		t.finalize()
	} else {
		t.blackHole.Wait()
	}

	assertValueIsWHNF("Thunk.result", t.result)

	return t.result
}

func (t *Thunk) lock(s thunkState) bool {
	for {
		switch t.loadState() {
		case normal:
			return false
		case app:
			if t.compareAndSwapState(app, s) {
				return true
			}
		}
	}
}

func (t *Thunk) delegateEval(parent *Thunk) (*Thunk, Arguments, bool) {
	if t.lock(spinLock) {
		f := t.swapFunction(identity)
		args := t.swapArguments(Arguments{[]*Thunk{parent}, nil, nil, nil})
		t.storeState(app)
		return f, args, true
	}

	return nil, Arguments{}, false
}

func (t *Thunk) swapFunction(new *Thunk) *Thunk {
	old := t.function
	t.function = new
	return old
}

func (t *Thunk) swapArguments(new Arguments) Arguments {
	old := t.args
	t.args = new
	return old
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

func (t *Thunk) loadState() thunkState {
	return thunkState(atomic.LoadInt32((*int32)(&t.state)))
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
