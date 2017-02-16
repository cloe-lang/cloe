package vm

import (
	"github.com/stretchr/testify/assert"
	"math"
	"testing"
)

func TestCompile(t *testing.T) {
	const n1, n2, n3 float64 = 2, 3, 4

	f := Compile(
		NewSignature(
			[]string{"f", "x1", "x2", "x3"}, []OptionalArgument{}, "",
			[]string{}, []OptionalArgument{}, "",
		),
		IRApp(0, NewIRPositionalArguments(1, IRApp(0, NewIRPositionalArguments(2, 3)))))

	x1 := float64(PApp(f, Pow, NewNumber(n1), NewNumber(n2), NewNumber(n3)).Eval().(NumberType))
	x2 := math.Pow(n1, math.Pow(n2, n3))

	t.Logf("%f == %f?", x1, x2)

	assert.Equal(t, x1, x2)
}
