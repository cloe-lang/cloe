package desugar

import (
	"github.com/cloe-lang/cloe/src/lib/ast"
	"github.com/cloe-lang/cloe/src/lib/debug"
	"github.com/cloe-lang/cloe/src/lib/gensym"
)

func desugarAnonymousFunctions(x any) []any {
	c := newAnonymousFunctionConverter()
	x = c.convert(x)
	return append(c.lets, x)
}

type anonymousFunctionConverter struct {
	lets []any
}

func newAnonymousFunctionConverter() anonymousFunctionConverter {
	return anonymousFunctionConverter{}
}

func (c *anonymousFunctionConverter) convert(x any) any {
	return ast.Convert(func(x any) any {
		switch x := x.(type) {
		case ast.AnonymousFunction:
			n := gensym.GenSym()
			c.lets = append(
				c.lets,
				desugarAnonymousFunctionsInDefFunction(
					ast.NewDefFunction(n, x.Signature(), x.Lets(), x.Body(), debug.NewGoInfo(0))))
			return n
		case ast.DefFunction:
			return desugarAnonymousFunctionsInDefFunction(x)
		}

		return nil
	}, x)
}

func desugarAnonymousFunctionsInDefFunction(f ast.DefFunction) ast.DefFunction {
	ls := make([]any, 0, len(f.Lets()))

	for _, l := range f.Lets() {
		ls = append(ls, desugarAnonymousFunctions(l)...)
	}

	c := newAnonymousFunctionConverter()
	b := c.convert(f.Body())
	ls = append(ls, c.lets...)

	return ast.NewDefFunction(f.Name(), f.Signature(), ls, b, f.DebugInfo())
}
