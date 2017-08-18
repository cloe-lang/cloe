package desugar

import (
	"github.com/tisp-lang/tisp/src/lib/ast"
	"github.com/tisp-lang/tisp/src/lib/debug"
	"github.com/tisp-lang/tisp/src/lib/gensym"
)

func desugarAnonymousFunctions(x interface{}) []interface{} {
	c := newAnonymousFunctionConverter()
	x = c.convert(x)
	return append(c.lets, x)
}

type anonymousFunctionConverter struct {
	lets []interface{}
}

func newAnonymousFunctionConverter() anonymousFunctionConverter {
	return anonymousFunctionConverter{}
}

func (c *anonymousFunctionConverter) convert(x interface{}) interface{} {
	return ast.Convert(func(x interface{}) interface{} {
		switch x := x.(type) {
		case ast.AnonymousFunction:
			n := gensym.GenSym()
			c.lets = append(
				c.lets,
				desugarAnonymousFunctionsInLetFunction(
					ast.NewLetFunction(n, x.Signature(), nil, x.Body(), debug.NewGoInfo(0))))
			return n
		case ast.LetFunction:
			return desugarAnonymousFunctionsInLetFunction(x)
		}

		return nil
	}, x)
}

func desugarAnonymousFunctionsInLetFunction(f ast.LetFunction) ast.LetFunction {
	ls := make([]interface{}, 0, len(f.Lets()))

	for _, l := range f.Lets() {
		ls = append(ls, desugarAnonymousFunctions(l)...)
	}

	c := newAnonymousFunctionConverter()
	b := c.convert(f.Body())
	ls = append(ls, c.lets...)

	return ast.NewLetFunction(f.Name(), f.Signature(), ls, b, f.DebugInfo())
}
