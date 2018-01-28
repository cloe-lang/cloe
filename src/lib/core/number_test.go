package core

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
)

var n1, n2 float64 = 123, 42

func TestNumberEqual(t *testing.T) {
	n := NewNumber(123)

	assert.True(t, testEqual(n, n))
	assert.True(t, !testEqual(n, NewNumber(456)))
}

func TestNumberAdd(t *testing.T) {
	assert.Equal(t, NewNumber(0), EvalPure(PApp(Add)).(*NumberType))
	assert.Equal(t,
		NewNumber(n1+n2),
		EvalPure(PApp(Add, NewNumber(n1), NewNumber(n2))).(*NumberType))
}

func TestNumberSub(t *testing.T) {
	assert.Equal(t,
		NewNumber(n1-n2),
		EvalPure(PApp(Sub, NewNumber(n1), NewNumber(n2))).(*NumberType))
}

func TestNumberMul(t *testing.T) {
	assert.Equal(t, NewNumber(1), EvalPure(PApp(Mul)).(*NumberType))
	assert.Equal(t,
		NewNumber(n1*n2),
		EvalPure(PApp(Mul, NewNumber(n1), NewNumber(n2))).(*NumberType))
}

func TestNumberDiv(t *testing.T) {
	assert.Equal(t,
		NewNumber(n1/n2),
		EvalPure(PApp(Div, NewNumber(n1), NewNumber(n2))).(*NumberType))
}

func TestNumberFloorDiv(t *testing.T) {
	assert.Equal(t,
		NewNumber(math.Floor(n1/n2)),
		EvalPure(PApp(FloorDiv, NewNumber(n1), NewNumber(n2))).(*NumberType))
}

func TestNumberMod(t *testing.T) {
	assert.Equal(t,
		NewNumber(math.Mod(float64(n1), float64(n2))),
		EvalPure(PApp(Mod, NewNumber(n1), NewNumber(n2))).(*NumberType))
}

func TestNumberPow(t *testing.T) {
	assert.Equal(t,
		NewNumber(math.Pow(float64(n1), float64(n2))),
		EvalPure(PApp(Pow, NewNumber(n1), NewNumber(n2))).(*NumberType))
}

func TestNumberFunctionsError(t *testing.T) {
	for _, v := range []Value{
		PApp(Add, Nil),
		App(Add, NewArguments([]PositionalArgument{NewPositionalArgument(Nil, true)}, nil, nil)),
		App(Add, NewArguments([]PositionalArgument{
			NewPositionalArgument(NewNumber(42), false),
			NewPositionalArgument(Nil, true),
		}, nil, nil)),
		PApp(Sub, Nil),
		App(Sub, NewArguments([]PositionalArgument{NewPositionalArgument(Nil, true)}, nil, nil)),
		App(Sub, NewArguments([]PositionalArgument{
			NewPositionalArgument(NewNumber(42), false),
			NewPositionalArgument(Nil, true),
		}, nil, nil)),
		App(Sub, NewArguments([]PositionalArgument{
			NewPositionalArgument(NewNumber(42), false),
			NewPositionalArgument(NewNumber(42), false),
			NewPositionalArgument(Nil, true),
		}, nil, nil)),
		App(Sub, NewArguments([]PositionalArgument{
			NewPositionalArgument(NewNumber(42), false),
			NewPositionalArgument(NewNumber(42), false),
			NewPositionalArgument(Nil, false),
		}, nil, nil)),
		PApp(Mod, Nil, NewNumber(42)),
		PApp(Mod, NewNumber(42), Nil),
	} {
		_, ok := EvalPure(v).(ErrorType)
		assert.True(t, ok)
	}
}

func TestNumberToString(t *testing.T) {
	for _, c := range []struct {
		expected string
		number   float64
	}{
		{"1", 1},
		{"1.1", 1.1},
	} {
		assert.Equal(t, NewString(c.expected), EvalPure(PApp(ToString, NewNumber(c.number))))
	}
}

func TestNumberCompare(t *testing.T) {
	for _, c := range []struct {
		answer      int
		left, right Value
	}{
		{0, NewNumber(42), NewNumber(42)},
		{1, NewNumber(1), NewNumber(0)},
		{1, NewNumber(0), NewNumber(-1)},
		{-1, NewNumber(0), NewNumber(1)},
		{-1, NewNumber(-1), NewNumber(0)},
	} {
		assert.Equal(t, c.answer, testCompare(c.left, c.right))
	}
}
