package core

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var impureFunction = NewLazyFunction(
	NewSignature([]string{"x"}, nil, "", nil, nil, ""),
	func(ts ...*Thunk) Value {
		return newEffect(ts[0])
	})

func TestThunkEvalWithNotCallable(t *testing.T) {
	assert.Equal(t, "TypeError", PApp(Nil).Eval().(ErrorType).name)
}

func TestThunkEvalWithImpureFunctionCall(t *testing.T) {
	assert.Equal(t, "ImpureFunctionError", PApp(impureFunction, Nil).Eval().(ErrorType).name)
}

func TestThunkEvalByCallingError(t *testing.T) {
	e := PApp(NewError("Apple", "pen.")).Eval().(ErrorType)
	t.Log(e)
	assert.Equal(t, 1, len(e.callTrace))
}

func TestThunkEvalByCallingErrorTwice(t *testing.T) {
	e := PApp(PApp(NewError("Apple", "pen."))).Eval().(ErrorType)
	t.Log(e)
	assert.Equal(t, 2, len(e.callTrace))
}

func TestThunkEvalEffect(t *testing.T) {
	s := NewString("foo")
	assert.Equal(t, s.Eval(), PApp(impureFunction, s).EvalEffect())
}

func TestThunkEvalEffectWithNonEffect(t *testing.T) {
	for _, th := range []*Thunk{Nil, PApp(identity, Nil)} {
		v := th.EvalEffect()
		err, ok := v.(ErrorType)
		t.Logf("%#v\n", v)
		assert.True(t, ok)
		assert.Equal(t, "TypeError", err.name)
	}
}

func TestThunkEvalEffectWithError(t *testing.T) {
	v := OutOfRangeError().EvalEffect()
	err, ok := v.(ErrorType)
	t.Logf("%#v\n", v)
	assert.True(t, ok)
	assert.Equal(t, "OutOfRangeError", err.name)
}

func TestNormal(t *testing.T) {
	for _, v := range []Value{Nil.Eval(), Nil} {
		_, ok := Normal(v).Eval().(ErrorType)
		assert.False(t, ok)
	}
}
