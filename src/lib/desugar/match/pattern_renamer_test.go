package match

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tisp-lang/tisp/src/lib/ast"
)

func TestPatternRenamerRenameFail(t *testing.T) {
	defer func() {
		assert.NotEqual(t, nil, recover())
	}()

	newPatternRenamer().rename(
		ast.NewSwitch("nil", []ast.SwitchCase{ast.NewSwitchCase("nil", "true")}, false))
}
