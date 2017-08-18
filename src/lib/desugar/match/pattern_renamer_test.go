package match

import (
	"testing"

	"github.com/tisp-lang/tisp/src/lib/ast"
)

func TestPatternRenamerRenameFail(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fail()
		}
	}()

	newPatternRenamer().rename(
		ast.NewSwitch("nil", []ast.SwitchCase{ast.NewSwitchCase("nil", "true")}, false))
}
