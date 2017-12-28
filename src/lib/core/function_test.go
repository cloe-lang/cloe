package core

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFunctionToString(t *testing.T) {
	assert.Equal(t, StringType("<function>"), PApp(ToString, If).Eval())
}

func TestNewEffectFunction(t *testing.T) {
	th := PApp(NewEffectFunction(
		NewSignature(nil, nil, "", nil, nil, ""),
		func(...*Thunk) Value { return Nil }))

	assert.Equal(t, "ImpureFunctionError", th.Eval().(ErrorType).Name())
	assert.Equal(t, Nil.Eval(), th.EvalEffect().(NilType))
}
