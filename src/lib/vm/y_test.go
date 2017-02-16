package vm

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestY(t *testing.T) {
	for _, n := range []float64{0, 1, 2, 3, 4, 5, 6, 42, 100} {
		n1 := lazyFactorial(NewNumber(n))
		n2 := strictFactorial(n)

		t.Logf("%d: %f == %f?\n", int(n), n1, n2)

		assert.Equal(t, n1, n2)
	}
}

func strictFactorial(n float64) float64 {
	if n == 0 {
		return 1
	}

	return n * strictFactorial(n-1)
}

func lazyFactorial(t *Thunk) float64 {
	return float64(PApp(PApp(Y, lazyFactorialImpl), t).Eval().(NumberType))
}

var lazyFactorialImpl = NewLazyFunction(
	NewSignature(
		[]string{"me", "num"}, []OptionalArgument{}, "",
		[]string{}, []OptionalArgument{}, "",
	),
	func(ts ...*Thunk) Object {
		return PApp(If,
			PApp(Equal, ts[1], NewNumber(0)),
			NewNumber(1),
			PApp(Mul,
				ts[1],
				PApp(ts[0], append([]*Thunk{PApp(Sub, ts[1], NewNumber(1))}, ts[2:]...)...)))
	})

func TestYsMultipleFs(t *testing.T) {
	evenWithExtraArg := NewLazyFunction(
		NewSignature(
			[]string{"fs", "dummyArg", "num"}, []OptionalArgument{}, "",
			[]string{}, []OptionalArgument{}, "",
		),
		func(ts ...*Thunk) Object {
			n := ts[2]

			return PApp(If,
				PApp(Equal, n, NewNumber(0)),
				True,
				PApp(PApp(First, PApp(Rest, ts[0])), PApp(Sub, n, NewNumber(1))))
		})

	odd := NewLazyFunction(
		NewSignature(
			[]string{"fs", "num"}, []OptionalArgument{}, "",
			[]string{}, []OptionalArgument{}, "",
		),
		func(ts ...*Thunk) Object {
			n := ts[1]

			return PApp(If,
				PApp(Equal, n, NewNumber(0)),
				False,
				PApp(PApp(First, ts[0]), Nil, PApp(Sub, n, NewNumber(1))))
		})

	fs := PApp(Ys, evenWithExtraArg, odd)

	for _, n := range []float64{0, 1, 2, 3, 4, 5, 6, 42, 100, 121, 256, 1023} {
		b1 := bool(PApp(PApp(First, fs), NewString("unused"), NewNumber(n)).Eval().(BoolType))
		b2 := bool(PApp(PApp(First, PApp(Rest, fs)), NewNumber(n)).Eval().(BoolType))

		t.Logf("n = %v, even? %v, odd? %v\n", n, b1, b2)

		rem := int(n) % 2
		assert.Equal(t, b1, rem == 0)
		assert.Equal(t, b2, rem != 0)
	}
}
