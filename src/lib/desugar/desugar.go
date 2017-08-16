package desugar

import "github.com/tisp-lang/tisp/src/lib/desugar/match"

// Desugar desugars a module of statements in AST.
func Desugar(ss []interface{}) []interface{} {
	for _, f := range []func(interface{}) []interface{}{
		desugarRecursiveLetVar,
		match.Desugar,
		desugarMutualRecursionStatement,
		desugarSelfRecursiveStatement,
		flattenStatement,
	} {
		new := make([]interface{}, 0, 2*len(ss))

		for _, s := range ss {
			new = append(new, f(s)...)
		}

		ss = new
	}

	return ss
}
