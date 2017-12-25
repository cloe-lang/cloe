package desugar

import (
	"github.com/coel-lang/coel/src/lib/ast"
)

func removeUnusedVariables(x interface{}) []interface{} {
	return []interface{}{ast.Convert(func(x interface{}) interface{} {
		f, ok := x.(ast.LetFunction)

		if !ok {
			return nil
		}

		ls := make([]interface{}, 0, len(f.Lets()))

		for _, l := range reverseSlice(f.Lets()) {
			for _, e := range append(letVarsToExpressions(ls), f.Body()) {
				if len(newNames(l.(ast.LetVar).Name()).findInExpression(e)) > 0 {
					ls = append(ls, l)
					break
				}
			}
		}

		return ast.NewLetFunction(f.Name(), f.Signature(), reverseSlice(ls), f.Body(), f.DebugInfo())
	}, x)}
}

func letVarsToExpressions(ls []interface{}) []interface{} {
	es := make([]interface{}, 0, len(ls))

	for _, l := range ls {
		es = append(es, l.(ast.LetVar).Expr())
	}

	return es
}
