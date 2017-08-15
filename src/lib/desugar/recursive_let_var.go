package desugar

import (
	"github.com/tisp-lang/tisp/src/lib/ast"
	"github.com/tisp-lang/tisp/src/lib/debug"
	"github.com/tisp-lang/tisp/src/lib/gensym"
)

func desugarRecursiveLetVar(s interface{}) []interface{} {
	l, ok := s.(ast.LetVar)

	if !ok || len(newNames(l.Name()).find(l.Expr())) == 0 {
		return []interface{}{s}
	}

	f := gensym.GenSym()

	return []interface{}{
		ast.NewLetFunction(
			l.Name(),
			ast.NewSignature(nil, nil, "", nil, nil, ""),
			nil,
			l.Expr(),
			debug.NewGoInfo(0)),
		ast.NewLetVar(f, l.Name()),
		ast.NewLetVar(l.Name(), ast.NewPApp(f, nil, debug.NewGoInfo(0)))}
}
