package desugar

import (
	"testing"

	"github.com/tisp-lang/tisp/src/lib/ast"
	"github.com/tisp-lang/tisp/src/lib/consts"
	"github.com/tisp-lang/tisp/src/lib/debug"
)

func TestDesugarDictionaryExpansion(t *testing.T) {
	desugarDictionaryExpansion(ast.NewLetVar(
		"foo",
		ast.NewApp(
			consts.Names.DictionaryFunction,
			ast.NewArguments([]ast.PositionalArgument{
				ast.NewPositionalArgument("foo", false),
				ast.NewPositionalArgument("bar", true),
			}, nil, nil),
			debug.NewGoInfo(0))))
}
