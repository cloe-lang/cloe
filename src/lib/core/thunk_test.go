package core

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var impureFunction = NewLazyFunction(
	NewSignature([]string{"x"}, nil, "", nil, nil, ""),
	func(vs ...Value) Value {
		return newEffect(vs[0])
	})

func TestThunkEvalWithNotCallable(t *testing.T) {
	assert.Equal(t, "TypeError", EvalPure(PApp(Nil)).(ErrorType).Name())
}

func TestThunkEvalWithImpureFunctionCall(t *testing.T) {
	assert.Equal(t, "ImpureFunctionError", EvalPure(PApp(impureFunction, Nil)).(ErrorType).Name())
}

func TestThunkEvalByCallingError(t *testing.T) {
	e := EvalPure(PApp(DummyError)).(ErrorType)
	t.Log(e)
	assert.Equal(t, 1, len(e.callTrace))
}

func TestThunkEvalByCallingErrorTwice(t *testing.T) {
	e := EvalPure(PApp(PApp(DummyError))).(ErrorType)
	t.Log(e)
	assert.Equal(t, 2, len(e.callTrace))
}

func TestThunkEvalImpure(t *testing.T) {
	s := NewString("foo")
	assert.Equal(t, s, EvalImpure(PApp(impureFunction, s)))
}

func TestThunkEvalImpureWithNonEffect(t *testing.T) {
	for _, v := range []Value{Nil, PApp(identity, Nil)} {
		v := EvalImpure(v)
		err, ok := v.(ErrorType)
		t.Logf("%#v\n", v)
		assert.True(t, ok)
		assert.Equal(t, "TypeError", err.Name())
	}
}

func TestThunkEvalImpureWithError(t *testing.T) {
	v := EvalImpure(DummyError)
	err, ok := v.(ErrorType)
	t.Logf("%#v\n", v)
	assert.True(t, ok)
	assert.Equal(t, "DummyError", err.Name())
}
