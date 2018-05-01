package desugar

import (
	"testing"

	"github.com/cloe-lang/cloe/src/lib/ast"
	"github.com/cloe-lang/cloe/src/lib/consts"
	"github.com/cloe-lang/cloe/src/lib/debug"
)

func TestDesugarDictionaryExpansion(t *testing.T) {
	desugarDictionaryExpansion(ast.NewLetVar(
		"foo",
		ast.NewApp(
			consts.Names.DictionaryFunction,
			ast.NewArguments([]ast.PositionalArgument{
				ast.NewPositionalArgument("foo", false),
				ast.NewPositionalArgument("bar", true),
			}, nil),
			debug.NewGoInfo(0))))
}
