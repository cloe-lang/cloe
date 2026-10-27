package desugar

import (
	"github.com/cloe-lang/cloe/src/lib/ast"
	"github.com/cloe-lang/cloe/src/lib/consts"
)

func desugarEmptyCollection(x any) []any {
	return []any{ast.Convert(func(x any) any {
		a, ok := x.(ast.App)

		if !ok {
			return nil
		}

		if len(a.Arguments().Positionals()) != 0 {
			return nil
		}

		switch a.Function() {
		case consts.Names.ListFunction:
			return consts.Names.EmptyList
		case consts.Names.DictionaryFunction:
			return consts.Names.EmptyDictionary
		}

		return nil
	}, x)}
}
