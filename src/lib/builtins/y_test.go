package builtins

import (
	"testing"

	"github.com/coel-lang/coel/src/lib/core"
	"github.com/coel-lang/coel/src/lib/systemt"
	"github.com/stretchr/testify/assert"
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

func lazyFactorial(t core.Value) float64 {
	return float64(*core.EvalPure(core.PApp(core.PApp(Y, lazyFactorialImpl), t)).(*core.NumberType))
}

var lazyFactorialImpl = core.NewLazyFunction(
	core.NewSignature([]string{"me", "num"}, nil, "", nil, nil, ""),
	func(ts ...core.Value) core.Value {
		return core.PApp(core.If,
			core.PApp(core.Equal, ts[1], core.NewNumber(0)),
			core.NewNumber(1),
			core.PApp(core.Mul,
				ts[1],
				core.PApp(ts[0], append([]core.Value{core.PApp(core.Sub, ts[1], core.NewNumber(1))}, ts[2:]...)...)))
	})

func BenchmarkYInfiniteRecursion(b *testing.B) {
	v := core.PApp(Y, core.NewLazyFunction(
		core.NewSignature([]string{"me"}, nil, "", nil, nil, ""),
		func(ts ...core.Value) core.Value {
			return ts[0]
		}))

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		v = core.PApp(core.EvalPure(v))
	}
}

func BenchmarkY(b *testing.B) {
	go systemt.RunDaemons()
	core.EvalPure(core.PApp(toZero, core.NewNumber(float64(b.N))))
}

func BenchmarkGoY(b *testing.B) {
	toZeroGo(float64(b.N))
}

var toZero = core.PApp(Y, core.NewLazyFunction(
	core.NewSignature([]string{"me", "num"}, nil, "", nil, nil, ""),
	func(ts ...core.Value) core.Value {
		return core.PApp(core.If,
			core.PApp(core.Equal, ts[1], core.NewNumber(0)),
			core.NewString("Benchmark finished!"),
			core.PApp(ts[0], core.PApp(core.Sub, ts[1], core.NewNumber(1))))
	}))

func toZeroGo(f float64) string {
	v := core.Value(core.NewNumber(f))
	n := *core.EvalPure(v).(*core.NumberType)

	for n > 0 {
		v = core.PApp(core.Sub, v, core.NewNumber(1))
		n = *core.EvalPure(v).(*core.NumberType)
	}

	return "Benchmark finished!"
}
