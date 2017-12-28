package core

import (
	"sync"
	"sync/atomic"

	"github.com/coel-lang/coel/src/lib/debug"
)

type thunkState int32

const (
	normal thunkState = iota
	app
	spinLock
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

// Normal creates a thunk of a normal value as its result.
func Normal(v Value) *Thunk {
	if t, ok := v.(*Thunk); ok {
		return t
	}

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

// evalAny evaluates a thunk and returns a pure or impure (effect) value.
func (t *Thunk) evalAny() Value {
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

			if ok := child.delegateEval(t); !ok {
				t.result = child.evalAny()
				t.chainError(t.result)
				break
			}
		}

		assertValueIsNormal("Thunk.result", t.result)

		t.finalize()
	} else {
		t.blackHole.Wait()
	}

	assertValueIsNormal("Thunk.result", t.result)

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

func (t *Thunk) delegateEval(parent *Thunk) bool {
	if t.lock(spinLock) {
		parent.function = t.swapFunction(identity)
		parent.args = t.swapArguments(Arguments{[]*Thunk{parent}, nil, nil, nil})
		parent.info = t.info
		t.storeState(app)
		return true
	}

	return false
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
		t.result = e.Chain(t.info)
		return true
	}

	return false
}

func assertValueIsNormal(s string, v Value) {
	if _, ok := v.(*Thunk); ok {
		panic(s + " is *Thunk")
	}
}

// Eval evaluates a pure value.
func (t *Thunk) Eval() Value {
	v := t.evalAny()

	if _, ok := v.(effectType); ok {
		return ImpureFunctionError(v).Eval().(ErrorType).Chain(t.info)
	}

	return t.result
}

// EvalEffect evaluates an effect expression.
func (t *Thunk) EvalEffect() Value {
	v := t.evalAny()
	e, ok := v.(effectType)

	if !ok {
		return NotEffectError(v).Eval().(ErrorType).Chain(t.info)
	}

	return e.value.Eval()
}
