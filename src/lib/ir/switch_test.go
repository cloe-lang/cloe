package ir

import (
	"testing"

	"github.com/coel-lang/coel/src/lib/core"
	"github.com/stretchr/testify/assert"
)

func TestNewSwitch(t *testing.T) {
	NewSwitch(0, []Case{NewCase(core.Nil, 1)}, core.Nil)
}

func TestNewSwitchNoDefaultCase(t *testing.T) {
	defer func() {
		assert.NotNil(t, recover())
	}()

	NewSwitch(0, []Case{}, nil)
}

func TestSwitchInFunction(t *testing.T) {
	f := CompileFunction(
		core.NewSignature([]string{"x"}, nil, "", nil, nil, ""),
		nil,
		NewSwitch(0, []Case{
			NewCase(core.NewString("foo"), core.NewNumber(42)),
			NewCase(core.True, core.NewNumber(1993)),
		}, core.NewNumber(2049)))

	assert.Equal(t, 42.0, float64(core.PApp(f, core.NewString("foo")).Eval().(core.NumberType)))
	assert.Equal(t, 1993.0, float64(core.PApp(f, core.True).Eval().(core.NumberType)))
	assert.Equal(t, 2049.0, float64(core.PApp(f, core.Nil).Eval().(core.NumberType)))
}
