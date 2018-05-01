package desugar

import (
	"github.com/cloe-lang/cloe/src/lib/ast"
	"github.com/cloe-lang/cloe/src/lib/gensym"
)

func desugarSelfRecursiveStatement(x interface{}) []interface{} {
	y := ast.Convert(func(x interface{}) interface{} {
		switch x := x.(type) {
		case ast.DefFunction:
			x = desugarInnerSelfRecursions(x)

			if len(newNames(x.Name()).findInDefFunction(x)) == 0 {
				return x
			}

			unrec := gensym.GenSym()

			return []interface{}{
				ast.NewDefFunction(
					unrec,
					prependPositionalsToSig([]string{x.Name()}, x.Signature()),
					x.Lets(),
					x.Body(),
					x.DebugInfo()),
				ast.NewLetVar(x.Name(), ast.NewPApp("$y", []interface{}{unrec}, x.DebugInfo())),
			}
		}

		return nil
	}, x)

	if ys, ok := y.([]interface{}); ok {
		return ys
	}

	return []interface{}{y}
}

func desugarInnerSelfRecursions(f ast.DefFunction) ast.DefFunction {
	ls := make([]interface{}, 0, 2*len(f.Lets()))

	for _, l := range f.Lets() {
		ls = append(ls, desugarSelfRecursiveStatement(l)...)
	}

	return ast.NewDefFunction(f.Name(), f.Signature(), ls, f.Body(), f.DebugInfo())
}
