package desugar

import (
	"github.com/coel-lang/coel/src/lib/desugar/match"
)

// Desugar desugars a module of statements in AST.
func Desugar(ss []interface{}) []interface{} {
	for _, f := range []func(interface{}) []interface{}{
		desugarLetMatch,
		desugarEmptyCollection,
		desugarDictionaryExpansion,
		desugarAnonymousFunctions,
		match.Desugar,
		desugarMutualRecursionStatement,
		desugarSelfRecursiveStatement,
		flattenStatement,
		removeUnusedVariables,
		removeAliases,
	} {
		new := make([]interface{}, 0, 2*len(ss))

		for _, s := range ss {
			new = append(new, f(s)...)
		}

		ss = new
	}

	return ss
}
