package vm

import (
	"github.com/stretchr/testify/assert"
	"math"
	"testing"
)

var n1, n2 float64 = 123, 42

func TestNumberEqual(t *testing.T) {
	n := NewNumber(123)

	assert.True(t, testEqual(n, n))
	assert.True(t, !testEqual(n, NewNumber(456)))
}

func TestNumberAdd(t *testing.T) {
	assert.Equal(t,
		float64(App(Add, NewNumber(n1), NewNumber(n2)).Eval().(numberType)),
		n1+n2)
}

func TestNumberSub(t *testing.T) {
	assert.Equal(t,
		float64(App(Sub, NewNumber(n1), NewNumber(n2)).Eval().(numberType)),
		n1-n2)
}

func TestNumberMul(t *testing.T) {
	assert.Equal(t,
		float64(App(Mul, NewNumber(n1), NewNumber(n2)).Eval().(numberType)),
		n1*n2)
}

func TestNumberDiv(t *testing.T) {
	assert.Equal(t,
		float64(App(Div, NewNumber(n1), NewNumber(n2)).Eval().(numberType)),
		n1/n2)
}

func TestNumberMod(t *testing.T) {
	assert.Equal(t,
		float64(App(Mod, NewNumber(n1), NewNumber(n2)).Eval().(numberType)),
		math.Mod(float64(n1), float64(n2)))
}

func TestNumberPow(t *testing.T) {
	assert.Equal(t,
		float64(App(Pow, NewNumber(n1), NewNumber(n2)).Eval().(numberType)),
		math.Pow(float64(n1), float64(n2)))
}
