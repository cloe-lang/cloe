package core

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var impureFunction = NewLazyFunction(
	NewSignature([]string{"x"}, nil, "", nil, nil, ""),
	func(ts ...*Thunk) Value {
		return NewOutput(ts[0])
	})

func TestThunkEvalWithNotCallable(t *testing.T) {
	assert.Equal(t, "TypeError", PApp(Nil).Eval().(ErrorType).name)
}

func TestThunkEvalWithImpureFunctionCall(t *testing.T) {
	assert.Equal(t, "TypeError", PApp(impureFunction, Nil).Eval().(ErrorType).name)
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

func TestThunkEvalOutput(t *testing.T) {
	s := NewString("foo")
	assert.Equal(t, s.Eval(), PApp(impureFunction, s).EvalOutput())
}

func TestThunkEvalOutputWithNonOutput(t *testing.T) {
	v := Nil.EvalOutput()
	err, ok := v.(ErrorType)
	t.Logf("%#v\n", v)
	assert.True(t, ok)
	assert.Equal(t, err.name, "TypeError")
}

func TestThunkEvalOutputWithError(t *testing.T) {
	v := OutOfRangeError().EvalOutput()
	err, ok := v.(ErrorType)
	t.Logf("%#v\n", v)
	assert.True(t, ok)
	assert.Equal(t, err.name, "OutOfRangeError")
}
