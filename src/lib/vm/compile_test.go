package vm

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCompile(t *testing.T) {
	var n1, n2, n3 float64 = 123, 456, 789

	f := Compile([]interface{}{0, 1, []interface{}{0, 2, 3}})

	x1 := float64(App(f, Add, NewNumber(n1), NewNumber(n2), NewNumber(n3)).Eval().(numberType))
	x2 := n1 + n2 + n3

	t.Logf("%f == %f?", x1, x2)

	assert.Equal(t, x1, x2)
}
