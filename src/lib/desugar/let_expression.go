package desugar

import (
	"github.com/cloe-lang/cloe/src/lib/ast"
	"github.com/cloe-lang/cloe/src/lib/debug"
)

func desugarLetExpression(x any) []any {
	return []any{convertLetExpression(x)}
}

func convertLetExpression(x any) any {
	return ast.Convert(func(x any) any {
		switch x := x.(type) {
		case ast.LetExpression:
			ls := x.Lets()
			y := convertLetExpression(x.Expr())

			for i := range ls {
				l := convertLetExpression(ls[len(ls)-1-i])

				switch l := l.(type) {
				case ast.LetVar:
					y = ast.NewPApp(
						ast.NewAnonymousFunction(ast.NewSignature([]string{l.Name()}, "", nil, ""), nil, y),
						[]any{l.Expr()},
						debug.NewGoInfo(0))
				case ast.LetMatch:
					y = ast.NewMatch(l.Expr(), []ast.MatchCase{ast.NewMatchCase(l.Pattern(), y)})
				default:
					panic("unreachable")
				}
			}

			return y
		}

		return nil
	}, x)
}
