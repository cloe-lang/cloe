package desugar

import (
	"github.com/coel-lang/coel/src/lib/ast"
	"github.com/coel-lang/coel/src/lib/consts"
)

func desugarEmptyCollection(x interface{}) []interface{} {
	return []interface{}{ast.Convert(func(x interface{}) interface{} {
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
