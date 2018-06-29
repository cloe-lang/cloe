package desugar

import (
	"fmt"

	"github.com/cloe-lang/cloe/src/lib/ast"
	"github.com/cloe-lang/cloe/src/lib/consts"
	"github.com/cloe-lang/cloe/src/lib/debug"
	"github.com/cloe-lang/cloe/src/lib/gensym"
	"github.com/cloe-lang/cloe/src/lib/scalar"
)

func desugarLetMatch(x interface{}) []interface{} {
	x = ast.Convert(convertLetMatch, x)

	if ls, ok := x.([]interface{}); ok {
		return ls
	}

	return []interface{}{x}
}

func convertLetMatch(x interface{}) interface{} {
	switch x := x.(type) {
	case ast.LetMatch:
		ns := []string{}

		ast.Convert(func(x interface{}) interface{} {
			n, ok := x.(string)

			if !ok {
				return nil
			}

			if n[:1] != "$" && !scalar.Defined(n) {
				ns = append(ns, n)
			}

			return nil
		}, x.Pattern())

		d := gensym.GenSym()

		ls := []interface{}{
			ast.NewLetVar(
				d,
				ast.NewMatch(
					x.Expr(),
					[]ast.MatchCase{
						ast.NewMatchCase(
							x.Pattern(),
							ast.NewPApp(
								consts.Names.DictionaryFunction,
								namesToPairsOfIndexAndName(ns),
								debug.NewGoInfo(0))),
					})),
		}

		for i, n := range ns {
			ls = append(ls,
				ast.NewLetVar(n, ast.NewPApp(
					consts.Names.IndexFunction,
					[]interface{}{d, fmt.Sprint(i)},
					debug.NewGoInfo(0))))
		}

		return ls
	case ast.DefFunction:
		return ast.NewDefFunction(
			x.Name(), x.Signature(), convertLets(x.Lets()), x.Body(), x.DebugInfo())
	case ast.LetExpression:
		return ast.NewLetExpression(convertLets(x.Lets()), x.Expr())
	}

	return nil
}

func namesToPairsOfIndexAndName(ns []string) []interface{} {
	args := make([]interface{}, 0, 2*len(ns))

	for i, n := range ns {
		args = append(args, fmt.Sprint(i), n)
	}

	return args
}

func convertLets(ls []interface{}) []interface{} {
	ns := make([]interface{}, 0, len(ls))

	for _, l := range ls {
		switch x := ast.Convert(convertLetMatch, l).(type) {
		case []interface{}:
			ns = append(ns, x...)
		default:
			ns = append(ns, x)
		}
	}

	return ns
}
