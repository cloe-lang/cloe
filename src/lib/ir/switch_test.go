package ir

import (
	"testing"

	"github.com/tisp-lang/tisp/src/lib/core"
)

func TestNewSwitch(t *testing.T) {
	NewSwitch([]core.Value{core.Nil.Eval()}, []int{1})
}

func TestNewSwitchUnmatchedNumbersOfPatternsAndValues(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fail()
		}
	}()

	NewSwitch(
		[]core.Value{core.NewNumber(123).Eval(), core.NewNumber(456).Eval()},
		[]int{1})
}

func TestNewSwitchNoPattern(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fail()
		}
	}()

	NewSwitch([]core.Value{}, []int{})
}

func TestSwitchCompileToDict(t *testing.T) {
	NewSwitch([]core.Value{core.Nil.Eval()}, []int{1}).compileToDict()
}
