package core

import (
	"testing"
	"time"

	"github.com/cloe-lang/cloe/src/lib/systemt"
	"github.com/stretchr/testify/assert"
)

func TestFunctionCallError(t *testing.T) {
	f := NewLazyFunction(
		NewSignature([]string{"foo"}, "", nil, ""),
		func(vs ...Value) Value { return vs[0] })

	for _, v := range []Value{
		App(f, NewArguments([]PositionalArgument{NewPositionalArgument(EmptyList, true)}, nil)),
		App(f, NewArguments(nil, []KeywordArgument{NewKeywordArgument("bar", Nil)})),
	} {
		_, ok := EvalPure(v).(*ErrorType)
		assert.True(t, ok)
	}
}

func TestFunctionToString(t *testing.T) {
	assert.Equal(t, NewString("<function>"), EvalPure(PApp(ToString, If)))
}

func TestStrictFunctionParallelization(t *testing.T) {
	go systemt.RunDaemons()

	f := NewStrictFunction(
		NewSignature([]string{"foo"}, "", []OptionalParameter{NewOptionalParameter("bar", Nil)}, ""),
		func(vs ...Value) Value { return vs[0] })

	EvalPure(App(f, NewArguments(nil, []KeywordArgument{NewKeywordArgument("bar", Nil)})))

	time.Sleep(100 * time.Millisecond)
}

func TestNewEffectFunction(t *testing.T) {
	v := PApp(NewEffectFunction(
		NewSignature(nil, "", nil, ""),
		func(...Value) Value { return Nil }))

	assert.Equal(t, Nil, EvalPure(v))
	assert.Equal(t, Nil, EvalImpure(v))
}

func TestPartial(t *testing.T) {
	ifFunc := func(vs ...Value) bool {
		return bool(*EvalPure(PApp(PApp(Partial, If, False, True), vs...)).(*BooleanType))
	}

	assert.True(t, ifFunc(True))
	assert.True(t, !ifFunc(False))
}

func TestPartialError(t *testing.T) {
	for _, a := range []Arguments{
		NewPositionalArguments(Nil),
		NewPositionalArguments(Prepend),
	} {
		_, ok := EvalPure(PApp(App(Partial, a))).(*ErrorType)
		assert.True(t, ok)
	}
}

func TestClosureToString(t *testing.T) {
	assert.Equal(t, NewString("<function>"), EvalPure(PApp(ToString, PApp(Partial, If, True))))
}
