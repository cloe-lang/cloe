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
				ast.NewLetVar(n, ast.NewPApp(d, []interface{}{fmt.Sprint(i)}, debug.NewGoInfo(0))))
		}

		return ls
	case ast.DefFunction:
		ls := make([]interface{}, 0, len(x.Lets()))

		for _, l := range x.Lets() {
			l = ast.Convert(convertLetMatch, l)

			switch x := l.(type) {
			case []interface{}:
				ls = append(ls, x...)
			default:
				ls = append(ls, x)
			}
		}

		return ast.NewDefFunction(x.Name(), x.Signature(), ls, x.Body(), x.DebugInfo())
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
