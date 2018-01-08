package core

import (
	"testing"
	"time"

	"github.com/coel-lang/coel/src/lib/systemt"
	"github.com/stretchr/testify/assert"
)

func TestFunctionCallError(t *testing.T) {
	f := NewLazyFunction(
		NewSignature([]string{"foo"}, nil, "", []string{"bar"}, nil, ""),
		func(ts ...*Thunk) Value { return ts[0] })

	for _, th := range []*Thunk{
		PApp(f, Nil),
		App(f, NewArguments(nil, []KeywordArgument{NewKeywordArgument("bar", Nil)}, nil)),
	} {
		_, ok := th.Eval().(ErrorType)
		assert.True(t, ok)
	}
}

func TestFunctionToString(t *testing.T) {
	assert.Equal(t, StringType("<function>"), PApp(ToString, If).Eval())
}

func TestStrictFunctionParallelization(t *testing.T) {
	go systemt.RunDaemons()

	f := NewStrictFunction(
		NewSignature([]string{"foo"}, nil, "", []string{"bar"}, nil, ""),
		func(ts ...*Thunk) Value { return ts[0] })

	App(f, NewArguments(nil, []KeywordArgument{NewKeywordArgument("bar", Nil)}, nil)).Eval()

	time.Sleep(100 * time.Millisecond)
}

func TestNewEffectFunction(t *testing.T) {
	th := PApp(NewEffectFunction(
		NewSignature(nil, nil, "", nil, nil, ""),
		func(...*Thunk) Value { return Nil }))

	assert.Equal(t, "ImpureFunctionError", th.Eval().(ErrorType).Name())
	assert.Equal(t, Nil.Eval(), th.EvalEffect().(NilType))
}
