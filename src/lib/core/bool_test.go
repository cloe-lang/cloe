package core

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBoolEqual(t *testing.T) {
	for _, vs := range [][2]Value{
		{True, True},
		{False, False},
	} {
		assert.True(t, testEqual(vs[0], vs[1]))
	}

	for _, vs := range [][2]Value{
		{True, False},
		{False, True},
	} {
		assert.True(t, !testEqual(vs[0], vs[1]))
	}
}

func TestBoolToString(t *testing.T) {
	test := func(s string, b bool) {
		assert.Equal(t, NewString(s), EvalPure(PApp(ToString, NewBool(b))))
	}

	test("true", true)
	test("false", false)
}

func TestIf(t *testing.T) {
	for _, vs := range [][]Value{
		{Nil},
		{True, Nil, False},
		{False, False, Nil},
		{False, False, True, Nil, False, True, False},
	} {
		assert.Equal(t, Nil, EvalPure(PApp(If, vs...)))
	}
}

func TestIfWithInvalidArguments(t *testing.T) {
	for _, a := range []Arguments{
		NewPositionalArguments(),
		NewPositionalArguments(Nil, Nil, Nil),
		NewPositionalArguments(False, Nil),
		NewArguments([]PositionalArgument{NewPositionalArgument(Nil, true)}, nil, nil),
		NewArguments(
			[]PositionalArgument{
				NewPositionalArgument(PApp(Prepend, True, NewError("FooError", "Hi!")), true),
			},
			nil,
			nil),
	} {
		_, ok := EvalPure(App(If, a)).(ErrorType)
		assert.True(t, ok)
	}
}

func BenchmarkIf(b *testing.B) {
	for i := 0; i < b.N; i++ {
		EvalPure(PApp(If, False, False, False, False, False, False, False, False, True))
	}
}
