package core

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBoolEqual(t *testing.T) {
	for _, ts := range [][2]*Thunk{
		{True, True},
		{False, False},
	} {
		assert.True(t, testEqual(ts[0], ts[1]))
	}

	for _, ts := range [][2]*Thunk{
		{True, False},
		{False, True},
	} {
		assert.True(t, !testEqual(ts[0], ts[1]))
	}
}

func TestBoolToString(t *testing.T) {
	test := func(s string, b bool) {
		assert.Equal(t, StringType(s), PApp(ToString, NewBool(b)).Eval())
	}

	test("true", true)
	test("false", false)
}

func TestIf(t *testing.T) {
	assert.Equal(t, Nil.Eval(), PApp(If, True, Nil, False).Eval())
	assert.Equal(t, Nil.Eval(), PApp(If, False, False, Nil).Eval())
}

func TestIfWithMultipleConditions(t *testing.T) {
	assert.Equal(t, Nil.Eval(), PApp(If, False, False, True, Nil, False, True, False).Eval())
}

func TestIfWithInvalidArguments(t *testing.T) {
	for _, th := range []*Thunk{
		PApp(If, Nil, Nil, Nil),
		PApp(If, True, Nil),
	} {
		_, ok := th.Eval().(ErrorType)
		assert.True(t, ok)
	}
}

func TestIfWithInvalidRestArgument(t *testing.T) {
	for _, a := range []Arguments{
		NewArguments([]PositionalArgument{NewPositionalArgument(Nil, true)}, nil, nil),
		NewArguments(
			[]PositionalArgument{
				NewPositionalArgument(PApp(Prepend, True, NewError("FooError", "Hi!")), true),
			},
			nil,
			nil),
	} {
		_, ok := App(If, a).Eval().(ErrorType)
		assert.True(t, ok)
	}
}
