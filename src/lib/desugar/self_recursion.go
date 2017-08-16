package desugar

import (
	"github.com/tisp-lang/tisp/src/lib/ast"
	"github.com/tisp-lang/tisp/src/lib/gensym"
)

func desugarSelfRecursiveStatement(x interface{}) []interface{} {
	y := ast.Convert(func(x interface{}) interface{} {
		switch x := x.(type) {
		case ast.LetFunction:
			x = desugarInnerSelfRecursions(x)

			if len(newNames(x.Name()).findInFunction(x)) == 0 {
				return x
			}

			unrec := gensym.GenSym(x.Name(), "unrec")

			return []interface{}{
				ast.NewLetFunction(
					unrec,
					prependPosReqsToSig([]string{x.Name()}, x.Signature()),
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

func desugarInnerSelfRecursions(f ast.LetFunction) ast.LetFunction {
	ls := make([]interface{}, 0, 2*len(f.Lets()))

	for _, l := range f.Lets() {
		ls = append(ls, desugarSelfRecursiveStatement(l)...)
	}

	return ast.NewLetFunction(f.Name(), f.Signature(), ls, f.Body(), f.DebugInfo())
}
