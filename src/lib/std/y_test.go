package std

import (
	"testing"

	"github.com/raviqqe/tisp/src/lib/core"
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

func lazyFactorial(t *core.Thunk) float64 {
	return float64(core.PApp(core.PApp(Y, lazyFactorialImpl), t).Eval().(core.NumberType))
}

var lazyFactorialImpl = core.NewLazyFunction(
	core.NewSignature(
		[]string{"me", "num"}, []core.OptionalArgument{}, "",
		[]string{}, []core.OptionalArgument{}, "",
	),
	func(ts ...*core.Thunk) core.Value {
		return core.PApp(core.If,
			core.PApp(core.Equal, ts[1], core.NewNumber(0)),
			core.NewNumber(1),
			core.PApp(core.Mul,
				ts[1],
				core.PApp(ts[0], append([]*core.Thunk{core.PApp(core.Sub, ts[1], core.NewNumber(1))}, ts[2:]...)...)))
	})

func TestYsMultipleFs(t *testing.T) {
	evenWithExtraArg := core.NewLazyFunction(
		core.NewSignature(
			[]string{"fs", "dummyArg", "num"}, []core.OptionalArgument{}, "",
			[]string{}, []core.OptionalArgument{}, "",
		),
		func(ts ...*core.Thunk) core.Value {
			n := ts[2]

			return core.PApp(core.If,
				core.PApp(core.Equal, n, core.NewNumber(0)),
				core.True,
				core.PApp(core.PApp(core.First, core.PApp(core.Rest, ts[0])), core.PApp(core.Sub, n, core.NewNumber(1))))
		})

	odd := core.NewLazyFunction(
		core.NewSignature(
			[]string{"fs", "num"}, []core.OptionalArgument{}, "",
			[]string{}, []core.OptionalArgument{}, "",
		),
		func(ts ...*core.Thunk) core.Value {
			n := ts[1]

			return core.PApp(core.If,
				core.PApp(core.Equal, n, core.NewNumber(0)),
				core.False,
				core.PApp(core.PApp(core.First, ts[0]), core.Nil, core.PApp(core.Sub, n, core.NewNumber(1))))
		})

	fs := core.PApp(Ys, evenWithExtraArg, odd)

	for _, n := range []float64{0, 1, 2, 3, 4, 5, 6, 42, 100, 121, 256, 1023} {
		b1 := bool(core.PApp(core.PApp(core.First, fs), core.NewString("unused"), core.NewNumber(n)).Eval().(core.BoolType))
		b2 := bool(core.PApp(core.PApp(core.First, core.PApp(core.Rest, fs)), core.NewNumber(n)).Eval().(core.BoolType))

		t.Logf("n = %v, even? %v, odd? %v\n", n, b1, b2)

		rem := int(n) % 2
		assert.Equal(t, b1, rem == 0)
		assert.Equal(t, b2, rem != 0)
	}
}
