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
		func(vs ...Value) Value { return vs[0] })

	for _, v := range []Value{
		PApp(f, Nil),
		App(f, NewArguments(nil, []KeywordArgument{NewKeywordArgument("bar", Nil)}, nil)),
	} {
		_, ok := EvalPure(v).(ErrorType)
		assert.True(t, ok)
	}
}

func TestFunctionToString(t *testing.T) {
	assert.Equal(t, NewString("<function>"), EvalPure(PApp(ToString, If)))
}

func TestStrictFunctionParallelization(t *testing.T) {
	go systemt.RunDaemons()

	f := NewStrictFunction(
		NewSignature([]string{"foo"}, nil, "", []string{"bar"}, nil, ""),
		func(vs ...Value) Value { return vs[0] })

	EvalPure(App(f, NewArguments(nil, []KeywordArgument{NewKeywordArgument("bar", Nil)}, nil)))

	time.Sleep(100 * time.Millisecond)
}

func TestNewEffectFunction(t *testing.T) {
	v := PApp(NewEffectFunction(
		NewSignature(nil, nil, "", nil, nil, ""),
		func(...Value) Value { return Nil }))

	assert.Equal(t, "ImpureFunctionError", EvalPure(v).(ErrorType).Name())
	assert.Equal(t, Nil, EvalImpure(v).(NilType))
}

func TestPartial(t *testing.T) {
	ifFunc := func(vs ...Value) bool {
		return bool(*EvalPure(PApp(PApp(Partial, If, False, True), vs...)).(*BoolType))
	}

	assert.True(t, ifFunc(True))
	assert.True(t, !ifFunc(False))
}

func TestPartialError(t *testing.T) {
	for _, a := range []Arguments{
		NewPositionalArguments(Nil),
		NewPositionalArguments(Prepend),
	} {
		_, ok := EvalPure(PApp(App(Partial, a))).(ErrorType)
		assert.True(t, ok)
	}
}

func TestClosureToString(t *testing.T) {
	assert.Equal(t, NewString("<function>"), EvalPure(PApp(ToString, PApp(Partial, If, True))))
}
