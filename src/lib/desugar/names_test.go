package desugar

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tisp-lang/tisp/src/lib/ast"
	"github.com/tisp-lang/tisp/src/lib/debug"
)

func TestNamesFindWithLetVar(t *testing.T) {
	n := "x"
	assert.True(t, newNames(n).find(ast.NewLetVar(n, n)).include(n))
}

func TestNamesFindWithLetFunction(t *testing.T) {
	n := "x"

	for _, test := range []struct {
		letFunc ast.LetFunction
		answer  bool
	}{
		{
			ast.NewLetFunction(
				n,
				ast.NewSignature(nil, nil, "", nil, nil, ""),
				nil,
				n,
				debug.NewGoInfo(0)),
			true,
		},
		{
			ast.NewLetFunction(
				n,
				ast.NewSignature(nil, nil, "", nil, nil, ""),
				[]interface{}{ast.NewLetVar(n, n)},
				n,
				debug.NewGoInfo(0)),
			true,
		},
		{
			ast.NewLetFunction(
				n,
				ast.NewSignature(nil, nil, "", nil, nil, ""),
				[]interface{}{ast.NewLetVar(n, "y")},
				n,
				debug.NewGoInfo(0)),
			false,
		},
	} {
		assert.Equal(t, test.answer, newNames(n).find(test.letFunc).include(n))
	}
}
