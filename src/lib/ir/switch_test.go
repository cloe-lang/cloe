package ir

import (
	"testing"

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
