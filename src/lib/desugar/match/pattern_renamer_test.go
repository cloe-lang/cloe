package match

import (
	"testing"

	"github.com/cloe-lang/cloe/src/lib/ast"
	"github.com/stretchr/testify/assert"
)

func TestPatternRenamerRenamePanic(t *testing.T) {
	assert.Panics(t, func() {
		newPatternRenamer().Rename(
			ast.NewSwitch("nil", []ast.SwitchCase{ast.NewSwitchCase("nil", "true")}, false))
	})
}
