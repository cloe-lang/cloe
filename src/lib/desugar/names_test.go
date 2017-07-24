package desugar

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tisp-lang/tisp/src/lib/ast"
)

func TestNamesFind(t *testing.T) {
	n := "x"
	assert.True(t, newNames(n).find(ast.NewLetVar(n, n)).include(n))
}
