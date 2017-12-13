package builtins

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/coel-lang/coel/src/lib/core"
)

func TestY(t *testing.T) {
	for _, n := range []float64{0, 1, 2, 3, 4, 5, 6, 42, 100} {
		n1 := lazyFactorial(core.NewNumber(n))
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

func lazyFactorial(t *core.Thunk) float64 {
	return float64(core.PApp(core.PApp(Y, lazyFactorialImpl), t).Eval().(core.NumberType))
}

var lazyFactorialImpl = core.NewLazyFunction(
	core.NewSignature(
		[]string{"me", "num"}, nil, "",
		nil, nil, "",
	),
	func(ts ...*core.Thunk) core.Value {
		return core.PApp(core.If,
			core.PApp(core.Equal, ts[1], core.NewNumber(0)),
			core.NewNumber(1),
			core.PApp(core.Mul,
				ts[1],
				core.PApp(ts[0], append([]*core.Thunk{core.PApp(core.Sub, ts[1], core.NewNumber(1))}, ts[2:]...)...)))
	})

func BenchmarkY(b *testing.B) {
	b.Log(core.PApp(toZero, core.NewNumber(10000)).Eval())
}

var toZero = core.PApp(Y, core.NewLazyFunction(
	core.NewSignature(
		[]string{"me", "num"}, nil, "",
		nil, nil, "",
	),
	func(ts ...*core.Thunk) core.Value {
		return core.PApp(core.If,
			core.PApp(core.Equal, ts[1], core.NewNumber(0)),
			core.NewString("Benchmark finished!"),

			core.PApp(ts[0], core.PApp(core.Sub, ts[1], core.NewNumber(1))))
	}))
