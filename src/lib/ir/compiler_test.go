package ir

import (
	"../vm"
	"github.com/stretchr/testify/assert"
	"math"
	"testing"
)

func TestCompileFunction(t *testing.T) {
	const n1, n2, n3 float64 = 2, 3, 4

	f := compileFunction(
		vm.NewSignature(
			[]string{"f", "x1", "x2", "x3"}, []vm.OptionalArgument{}, "",
			[]string{}, []vm.OptionalArgument{}, "",
		),
		IRApp(0, NewIRPositionalArguments(1, IRApp(0, NewIRPositionalArguments(2, 3)))))

	x1 := float64(vm.PApp(f, vm.Pow, vm.NewNumber(n1), vm.NewNumber(n2), vm.NewNumber(n3)).Eval().(vm.NumberType))
	x2 := math.Pow(n1, math.Pow(n2, n3))

	t.Logf("%f == %f?", x1, x2)

	assert.Equal(t, x1, x2)
}
