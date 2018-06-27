package builtins

import (
	"testing"

	"github.com/cloe-lang/cloe/src/lib/core"
	"github.com/cloe-lang/cloe/src/lib/systemt"
	"github.com/stretchr/testify/assert"
)

func TestPar(t *testing.T) {
	go systemt.RunDaemons()

	n := core.PApp(Par, core.True, core.False, core.DummyError, core.Nil)

	assert.True(t, bool(*core.EvalPure(core.PApp(core.Equal, core.Nil, n)).(*core.BooleanType)))
}

func TestParError(t *testing.T) {
	go systemt.RunDaemons()

	_, ok := core.EvalPure(core.PApp(Par,
		core.True,
		core.False,
		core.Nil,
		core.DummyError,
	)).(*core.ErrorType)

	assert.True(t, ok)
}

func TestParWithNoArgument(t *testing.T) {
	go systemt.RunDaemons()

	_, ok := core.EvalPure(core.PApp(Par)).(*core.ErrorType)
	assert.True(t, ok)
}
