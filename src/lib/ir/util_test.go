package ir

import (
	"math"
	"testing"

	"github.com/coel-lang/coel/src/lib/core"
	"github.com/coel-lang/coel/src/lib/debug"
	"github.com/stretchr/testify/assert"
)

func TestCompileFunction(t *testing.T) {
	const n1, n2, n3 = 2, 3, 4

	f := CompileFunction(
		core.NewSignature(
			[]string{"f", "x1", "x2", "x3"}, nil, "",
			nil, nil, "",
		),
		nil,
		newAppWithDummyInfo(0, newPositionalArguments(1, newAppWithDummyInfo(0, newPositionalArguments(2, 3)))))

	x1 := float64(core.PApp(f, core.Pow, core.NewNumber(n1), core.NewNumber(n2), core.NewNumber(n3)).Eval().(core.NumberType))
	x2 := math.Pow(n1, math.Pow(n2, n3))

	t.Logf("%f == %f?", x1, x2)
	assert.Equal(t, x1, x2)
}

func TestCompileFunctionWithVars(t *testing.T) {
	const n1, n2, n3 = 2, 3, 4

	f := CompileFunction(
		core.NewSignature(
			[]string{"f", "x1", "x2", "x3"}, nil, "",
			nil, nil, "",
		),
		[]interface{}{newAppWithDummyInfo(0, newPositionalArguments(2, 3))},
		newAppWithDummyInfo(0, newPositionalArguments(1, 4)))

	x1 := float64(core.PApp(f, core.Pow, core.NewNumber(n1), core.NewNumber(n2), core.NewNumber(n3)).Eval().(core.NumberType))
	x2 := math.Pow(n1, math.Pow(n2, n3))

	t.Logf("%f == %f?", x1, x2)
	assert.Equal(t, x1, x2)
}

func newPositionalArguments(xs ...interface{}) Arguments {
	ps := make([]PositionalArgument, 0, len(xs))

	for _, x := range xs {
		ps = append(ps, NewPositionalArgument(x, false))
	}

	return NewArguments(ps, nil, nil)
}

func newAppWithDummyInfo(f interface{}, args Arguments) App {
	return NewApp(f, args, debug.NewGoInfo(0))
}

func TestInterpretExpressionFail(t *testing.T) {
	defer func() {
		assert.NotEqual(t, nil, recover())
	}()

	interpretExpression(nil, "foo")
}
