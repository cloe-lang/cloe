package builtins

import (
	"testing"

	"github.com/cloe-lang/cloe/src/lib/core"
	"github.com/stretchr/testify/assert"
)

func TestYsMultipleFs(t *testing.T) {
	evenWithExtraArg := core.NewLazyFunction(
		core.NewSignature(
			[]string{"fs", "dummyArg", "num"}, "",
			nil, "",
		),
		func(ts ...core.Value) core.Value {
			n := ts[2]

			return core.PApp(core.If,
				core.PApp(core.Equal, n, core.NewNumber(0)),
				core.True,
				core.PApp(core.PApp(core.First, core.PApp(core.Rest, ts[0])), core.PApp(core.Sub, n, core.NewNumber(1))))
		})

	odd := core.NewLazyFunction(
		core.NewSignature(
			[]string{"fs", "num"}, "",
			nil, "",
		),
		func(ts ...core.Value) core.Value {
			n := ts[1]

			return core.PApp(core.If,
				core.PApp(core.Equal, n, core.NewNumber(0)),
				core.False,
				core.PApp(core.PApp(core.First, ts[0]), core.Nil, core.PApp(core.Sub, n, core.NewNumber(1))))
		})

	fs := core.PApp(Ys, evenWithExtraArg, odd)

	for _, n := range []float64{0, 1, 2, 3, 4, 5, 6, 42, 100, 121, 256, 1023} {
		b1 := bool(*core.EvalPure(core.PApp(core.PApp(core.First, fs), core.NewString("unused"), core.NewNumber(n))).(*core.BooleanType))
		b2 := bool(*core.EvalPure(core.PApp(core.PApp(core.First, core.PApp(core.Rest, fs)), core.NewNumber(n))).(*core.BooleanType))

		t.Logf("n = %v, even? %v, odd? %v\n", n, b1, b2)

		rem := int(n) % 2
		assert.Equal(t, b1, rem == 0)
		assert.Equal(t, b2, rem != 0)
	}
}

func TestYsWithErroneousArgument(t *testing.T) {
	v := core.EvalPure(core.App(
		Ys,
		core.NewArguments(
			[]core.PositionalArgument{core.NewPositionalArgument(core.DummyError, true)},
			nil)))
	_, ok := v.(*core.ErrorType)
	assert.True(t, ok)
}
