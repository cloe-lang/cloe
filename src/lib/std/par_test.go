package std

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tisp-lang/tisp/src/lib/core"
	"github.com/tisp-lang/tisp/src/lib/systemt"
)

func TestPar(t *testing.T) {
	go systemt.RunDaemons()

	n := core.PApp(Par, core.True, core.False, core.ValueError("I am the error."), core.Nil)

	assert.True(t, bool(core.PApp(core.Equal, core.Nil, n).Eval().(core.BoolType)))
}

func TestParError(t *testing.T) {
	go systemt.RunDaemons()

	_, ok := core.PApp(Par,
		core.True,
		core.False,
		core.Nil,
		core.ValueError("I am the error."),
	).Eval().(core.ErrorType)

	assert.True(t, ok)
}

func TestParWithNoArgument(t *testing.T) {
	go systemt.RunDaemons()

	_, ok := core.PApp(Par).Eval().(core.ErrorType)
	assert.True(t, ok)
}
