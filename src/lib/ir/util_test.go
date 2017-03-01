package ir

import (
	"github.com/raviqqe/tisp/src/lib/core"
	"github.com/stretchr/testify/assert"
	"math"
	"testing"
)

func TestCompileFunction(t *testing.T) {
	const n1, n2, n3 = 2, 3, 4

	f := CompileFunction(
		core.NewSignature(
			[]string{"f", "x1", "x2", "x3"}, []core.OptionalArgument{}, "",
			[]string{}, []core.OptionalArgument{}, "",
		),
		[]interface{}{},
		NewApp(0, newPositionalArguments(1, NewApp(0, newPositionalArguments(2, 3)))))

	x1 := float64(core.PApp(f, core.Pow, core.NewNumber(n1), core.NewNumber(n2), core.NewNumber(n3)).Eval().(core.NumberType))
	x2 := math.Pow(n1, math.Pow(n2, n3))

	t.Logf("%f == %f?", x1, x2)
	assert.Equal(t, x1, x2)
}

func TestCompileFunctionWithVars(t *testing.T) {
	const n1, n2, n3 = 2, 3, 4

	f := CompileFunction(
		core.NewSignature(
			[]string{"f", "x1", "x2", "x3"}, []core.OptionalArgument{}, "",
			[]string{}, []core.OptionalArgument{}, "",
		),
		[]interface{}{NewApp(0, newPositionalArguments(2, 3))},
		NewApp(0, newPositionalArguments(1, 4)))

	x1 := float64(core.PApp(f, core.Pow, core.NewNumber(n1), core.NewNumber(n2), core.NewNumber(n3)).Eval().(core.NumberType))
	x2 := math.Pow(n1, math.Pow(n2, n3))

	t.Logf("%f == %f?", x1, x2)
	assert.Equal(t, x1, x2)
}

func newPositionalArguments(xs ...interface{}) Arguments {
	ps := make([]PositionalArgument, len(xs))

	for i, x := range xs {
		ps[i] = NewPositionalArgument(x, false)
	}

	return NewArguments(ps, []KeywordArgument{}, []interface{}{})
}
