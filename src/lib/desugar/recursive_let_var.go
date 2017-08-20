package desugar

import (
	"github.com/tisp-lang/tisp/src/lib/ast"
	"github.com/tisp-lang/tisp/src/lib/debug"
	"github.com/tisp-lang/tisp/src/lib/gensym"
)

func desugarRecursiveLetVar(x interface{}) []interface{} {
	y := ast.Convert(convertRecursiveLetVar, x)

	if ys, ok := y.([]interface{}); ok {
		return ys
	}

	return []interface{}{y}
}

func convertRecursiveLetVar(x interface{}) interface{} {
	switch x := x.(type) {
	case ast.LetVar:
		if len(newNames(x.Name()).findInExpression(x.Expr())) == 0 {
			return x
		}

		n := x.Name()
		f := gensym.GenSym()

		return []interface{}{
			ast.NewLetFunction(
				n,
				ast.NewSignature(nil, nil, "", nil, nil, ""),
				nil,
				ast.Convert(func(x interface{}) interface{} {
					if m, ok := x.(string); ok && n == m {
						return ast.NewPApp(n, nil, debug.NewGoInfo(0))
					}

					return nil
				}, x.Expr()),
				debug.NewGoInfo(0)),
			ast.NewLetVar(f, n),
			ast.NewLetVar(n, ast.NewPApp(f, nil, debug.NewGoInfo(0)))}
	case ast.LetFunction:
		ls := make([]interface{}, 0, 3*len(x.Lets()))

		for _, l := range x.Lets() {
			ls = append(ls, desugarRecursiveLetVar(l)...)
		}

		return ast.NewLetFunction(x.Name(), x.Signature(), ls, x.Body(), x.DebugInfo())
	}

	return nil
}
