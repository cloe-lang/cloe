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
	assert.Equal(t, float64(PApp(Add).Eval().(NumberType)), float64(0))
	assert.Equal(t,
		float64(PApp(Add, NewNumber(n1), NewNumber(n2)).Eval().(NumberType)),
		n1+n2)
}

func TestNumberSub(t *testing.T) {
	assert.Equal(t,
		float64(PApp(Sub, NewNumber(n1), NewNumber(n2)).Eval().(NumberType)),
		n1-n2)
}

func TestNumberMul(t *testing.T) {
	assert.Equal(t, float64(PApp(Mul).Eval().(NumberType)), float64(1))
	assert.Equal(t,
		float64(PApp(Mul, NewNumber(n1), NewNumber(n2)).Eval().(NumberType)),
		n1*n2)
}

func TestNumberDiv(t *testing.T) {
	assert.Equal(t,
		float64(PApp(Div, NewNumber(n1), NewNumber(n2)).Eval().(NumberType)),
		n1/n2)
}

func TestNumberMod(t *testing.T) {
	assert.Equal(t,
		float64(PApp(Mod, NewNumber(n1), NewNumber(n2)).Eval().(NumberType)),
		math.Mod(float64(n1), float64(n2)))
}

func TestNumberPow(t *testing.T) {
	assert.Equal(t,
		float64(PApp(Pow, NewNumber(n1), NewNumber(n2)).Eval().(NumberType)),
		math.Pow(float64(n1), float64(n2)))
}

func TestNumberToString(t *testing.T) {
	for _, xs := range []struct {
		expected string
		number   float64
	}{
		{"1", 1},
		{"1.1", 1.1},
	} {
		assert.Equal(t, StringType(xs.expected), PApp(ToString, NewNumber(xs.number)).Eval())
	}
}
