package core

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestThunkEval1Fail(t *testing.T) {
	e := PApp(NewError("Apple", "pen.")).Eval().(ErrorType)
	t.Log(e)
	assert.Equal(t, 1, len(e.callTrace))
}

func TestThunkEval2Fail(t *testing.T) {
	e := PApp(PApp(NewError("Apple", "pen."))).Eval().(ErrorType)
	t.Log(e)
	assert.Equal(t, 2, len(e.callTrace))
}

func TestThunkEvalOutputFail(t *testing.T) {
	v := Nil.EvalOutput()
	_, ok := v.(ErrorType)
	t.Logf("%#v\n", v)
	assert.True(t, ok)
}
