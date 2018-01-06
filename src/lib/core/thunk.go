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
	info      *debug.Info
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
func AppWithInfo(f *Thunk, args Arguments, i *debug.Info) *Thunk {
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
			v := t.function.Eval()
			t.function = nil

			f, ok := v.(callable)

			if !ok {
				t.result = NotCallableError(v).Eval()
				t.args = Arguments{}
				break
			}

			t.result = f.call(t.args)
			t.args = Arguments{}

			child, ok := t.result.(*Thunk)

			if !ok {
				break
			}

			if !child.delegateEval(t) {
				t.result = child.evalAny()
				break
			}
		}

		if e, ok := t.result.(ErrorType); ok {
			t.result = e.Chain(t.info)
		}

		// No need to clean up or finalize t.function, t.args, and t.state here
		// because of invariants.
		t.blackHole.Done()
	} else {
		t.blackHole.Wait()
	}

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
		parent.function, t.function = t.function, identity
		parent.args, t.args = t.args, NewPositionalArguments(parent)
		parent.info = t.info
		t.storeState(app)
		return true
	}

	return false
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

// Eval evaluates a pure value.
func (t *Thunk) Eval() Value {
	if _, ok := t.evalAny().(effectType); ok {
		return impureFunctionError().Eval().(ErrorType).Chain(t.info)
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

// EvalBool evaluates a thunk which is expected to be a boolean value.
func (t *Thunk) EvalBool() (BoolType, Value) {
	v := t.Eval()
	b, ok := v.(BoolType)

	if !ok {
		return false, NotBoolError(v).Eval()
	}

	return b, nil
}

// EvalList evaluates a thunk which is expected to be a list value.
func (t *Thunk) EvalList() (ListType, Value) {
	v := t.Eval()
	l, ok := v.(ListType)

	if !ok {
		return emptyList, NotListError(v).Eval()
	}

	return l, nil
}

// EvalNumber evaluates a thunk which is expected to be a number value.
func (t *Thunk) EvalNumber() (NumberType, Value) {
	v := t.Eval()
	n, ok := v.(NumberType)

	if !ok {
		return 0, NotNumberError(v).Eval()
	}

	return n, nil
}

// EvalString evaluates a thunk which is expected to be a string value.
func (t *Thunk) EvalString() (StringType, Value) {
	v := t.Eval()
	s, ok := v.(StringType)

	if !ok {
		return "", NotStringError(v).Eval()
	}

	return s, nil
}
