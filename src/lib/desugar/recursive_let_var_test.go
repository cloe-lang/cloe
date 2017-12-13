package desugar

import (
	"testing"

	"github.com/coel-lang/coel/src/lib/ast"
	"github.com/coel-lang/coel/src/lib/debug"
	"github.com/stretchr/testify/assert"
)

func TestDesugarRecursiveLetVar(t *testing.T) {
	l := ast.NewLetVar("foo", ast.NewPApp("prepend", []interface{}{"42", "foo"}, debug.NewGoInfo(0)))

	assert.True(t, len(desugarRecursiveLetVar(l)) > 1)
}

func TestDesugarRecursiveLetVarWithLetFunction(t *testing.T) {
	l := ast.NewLetFunction(
		"foo",
		ast.NewSignature(nil, nil, "", nil, nil, ""),
		nil,
		"123",
		debug.NewGoInfo(0))

	assert.Equal(t, 1, len(desugarRecursiveLetVar(l)))
}
