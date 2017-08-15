package desugar

import (
	"github.com/tisp-lang/tisp/src/lib/ast"
	"github.com/tisp-lang/tisp/src/lib/debug"
	"github.com/tisp-lang/tisp/src/lib/gensym"
)

func desugarRecursiveLetVar(s interface{}) []interface{} {
	switch s := s.(type) {
	case ast.LetVar:
		if len(newNames(s.Name()).find(s.Expr())) == 0 {
			return []interface{}{s}
		}

		n := s.Name()
		f := gensym.GenSym()

		return []interface{}{
			ast.NewLetFunction(
				n,
				ast.NewSignature(nil, nil, "", nil, nil, ""),
				nil,
				s.Expr(),
				debug.NewGoInfo(0)),
			ast.NewLetVar(f, n),
			ast.NewLetVar(n, ast.NewPApp(f, nil, debug.NewGoInfo(0)))}
	case ast.LetFunction:
		ls := make([]interface{}, 0, 3*len(s.Lets()))

		for _, l := range s.Lets() {
			ls = append(ls, desugarRecursiveLetVar(l)...)
		}

		return []interface{}{ast.NewLetFunction(s.Name(), s.Signature(), ls, s.Body(), s.DebugInfo())}
	case ast.MutualRecursion:
		fs := make([]ast.LetFunction, 0, len(s.LetFunctions()))

		for _, f := range s.LetFunctions() {
			for _, l := range desugarRecursiveLetVar(f) {
				fs = append(fs, l.(ast.LetFunction))
			}
		}

		return []interface{}{ast.NewMutualRecursion(fs, s.DebugInfo())}
	}

	return []interface{}{s}
}
