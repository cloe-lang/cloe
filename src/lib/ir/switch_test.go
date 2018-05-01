package ir

import (
	"testing"

	"github.com/cloe-lang/cloe/src/lib/core"
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
		core.NewSignature([]string{"x"}, "", nil, ""),
		nil,
		NewSwitch(0, []Case{
			NewCase(core.NewString("foo"), core.NewNumber(42)),
			NewCase(core.True, core.NewNumber(1993)),
		}, core.NewNumber(2049)))

	for _, c := range []struct{ answer, argument core.Value }{
		{core.NewNumber(42), core.NewString("foo")},
		{core.NewNumber(1993), core.True},
		{core.NewNumber(2049), core.Nil},
	} {
		assert.Equal(t, c.answer, core.EvalPure(core.PApp(f, c.argument)))
	}
}
