package ir

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tisp-lang/tisp/src/lib/core"
)

func TestNewSwitch(t *testing.T) {
	NewSwitch(0, []Case{NewCase(core.Nil, 1)})
}

func TestNewSwitchNoPattern(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fail()
		}
	}()

	NewSwitch(0, []Case{})
}

func TestSwitchCompileToDict(t *testing.T) {
	NewSwitch(0, []Case{NewCase(core.Nil, 1)}).compileToDict()
}

func TestSwitchInFunction(t *testing.T) {
	f := CompileFunction(
		core.NewSignature([]string{"x"}, nil, "", nil, nil, ""),
		nil,
		NewSwitch(0, []Case{
			NewCase(core.NewString("foo"), core.NewNumber(42)),
			NewCase(core.True, core.NewNumber(2049)),
		}))

	assert.Equal(t, 42.0, float64(core.PApp(f, core.NewString("foo")).Eval().(core.NumberType)))
	assert.Equal(t, 2049.0, float64(core.PApp(f, core.True).Eval().(core.NumberType)))
}
