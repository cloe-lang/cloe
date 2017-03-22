package desugar

import (
	"github.com/raviqqe/tisp/src/lib/ast"
	"github.com/raviqqe/tisp/src/lib/gensym"
)

func desugarSelfRecursiveStatement(s interface{}) []interface{} {
	switch s := s.(type) {
	case ast.LetFunction:
		return desugarSelfRecursiveFunction(s)
	default:
		return []interface{}{s}
	}
}

func desugarSelfRecursiveFunction(f ast.LetFunction) []interface{} {
	f = desugarInnerSelfRecursiveStatements(f)

	if ns := newNames(f.Name()); signatureToNames(f.Signature()).include(f.Name()) || len(ns.find(f.Lets()))+len(ns.find(f.Body())) == 0 {
		return []interface{}{f}
	}

	unrecursive := gensym.GenSym(f.Name(), "unrecursive")

	return []interface{}{
		ast.NewLetFunction(
			unrecursive,
			prependPosReqsToSig(f.Signature(), []string{f.Name()}),
			f.Lets(),
			f.Body(),
			f.DebugInfo()),
		ast.NewLetVar(
			f.Name(),
			ast.NewApp(
				"$y",
				ast.NewArguments([]ast.PositionalArgument{ast.NewPositionalArgument(unrecursive, false)}, []ast.KeywordArgument{}, []interface{}{}),
				f.DebugInfo())),
	}
}

func desugarInnerSelfRecursiveStatements(f ast.LetFunction) ast.LetFunction {
	ls := make([]interface{}, 0, 2*len(f.Lets()))

	for _, l := range f.Lets() {
		ls = append(ls, desugarSelfRecursiveStatement(l)...)
	}

	return ast.NewLetFunction(f.Name(), f.Signature(), ls, f.Body(), f.DebugInfo())
}
