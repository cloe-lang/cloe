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
	assert.Equal(t, NumberType(0), PApp(Add).Eval().(NumberType))
	assert.Equal(t,
		NumberType(n1+n2),
		PApp(Add, NewNumber(n1), NewNumber(n2)).Eval().(NumberType))
}

func TestNumberSub(t *testing.T) {
	assert.Equal(t,
		NumberType(n1-n2),
		PApp(Sub, NewNumber(n1), NewNumber(n2)).Eval().(NumberType))
}

func TestNumberMul(t *testing.T) {
	assert.Equal(t, NumberType(1), PApp(Mul).Eval().(NumberType))
	assert.Equal(t,
		NumberType(n1*n2),
		PApp(Mul, NewNumber(n1), NewNumber(n2)).Eval().(NumberType))
}

func TestNumberDiv(t *testing.T) {
	assert.Equal(t,
		NumberType(n1/n2),
		PApp(Div, NewNumber(n1), NewNumber(n2)).Eval().(NumberType))
}

func TestNumberFloorDiv(t *testing.T) {
	assert.Equal(t,
		NumberType(math.Floor(n1/n2)),
		PApp(FloorDiv, NewNumber(n1), NewNumber(n2)).Eval().(NumberType))
}

func TestNumberMod(t *testing.T) {
	assert.Equal(t,
		NumberType(math.Mod(float64(n1), float64(n2))),
		PApp(Mod, NewNumber(n1), NewNumber(n2)).Eval().(NumberType))
}

func TestNumberPow(t *testing.T) {
	assert.Equal(t,
		NumberType(math.Pow(float64(n1), float64(n2))),
		PApp(Pow, NewNumber(n1), NewNumber(n2)).Eval().(NumberType))
}

func TestNumberToString(t *testing.T) {
	for _, c := range []struct {
		expected string
		number   float64
	}{
		{"1", 1},
		{"1.1", 1.1},
	} {
		assert.Equal(t, StringType(c.expected), PApp(ToString, NewNumber(c.number)).Eval())
	}
}

func TestNumberCompare(t *testing.T) {
	assert.True(t, testCompare(NewNumber(42), NewNumber(42)) == 0)
	assert.True(t, testCompare(NewNumber(0), NewNumber(1)) == -1)
	assert.True(t, testCompare(NewNumber(1), NewNumber(0)) == 1)
	assert.True(t, testCompare(NewNumber(-1), NewNumber(0)) == -1)
	assert.True(t, testCompare(NewNumber(0), NewNumber(-1)) == 1)
}
