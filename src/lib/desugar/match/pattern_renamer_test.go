package match

import (
	"testing"

	"github.com/cloe-lang/cloe/src/lib/ast"
	"github.com/stretchr/testify/assert"
)

func TestPatternRenamerRenameFail(t *testing.T) {
	defer func() {
		assert.NotNil(t, recover())
	}()

	newPatternRenamer().rename(
		ast.NewSwitch("nil", []ast.SwitchCase{ast.NewSwitchCase("nil", "true")}, false))
}
