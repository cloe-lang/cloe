package builtins

import (
	"testing"

	"github.com/coel-lang/coel/src/lib/core"
	"github.com/coel-lang/coel/src/lib/systemt"
	"github.com/stretchr/testify/assert"
)

func TestPar(t *testing.T) {
	go systemt.RunDaemons()

	n := core.PApp(Par, core.True, core.False, core.ValueError("I am the error."), core.Nil)

	assert.True(t, bool(*core.EvalPure(core.PApp(core.Equal, core.Nil, n)).(*core.BoolType)))
}

func TestParError(t *testing.T) {
	go systemt.RunDaemons()

	_, ok := core.EvalPure(core.PApp(Par,
		core.True,
		core.False,
		core.Nil,
		core.ValueError("I am the error."),
	)).(core.ErrorType)

	assert.True(t, ok)
}

func TestParWithNoArgument(t *testing.T) {
	go systemt.RunDaemons()

	_, ok := core.EvalPure(core.PApp(Par)).(core.ErrorType)
	assert.True(t, ok)
}
