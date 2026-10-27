package desugar

import (
	"github.com/cloe-lang/cloe/src/lib/desugar/match"
)

// Desugar desugars a module of statements in AST.
func Desugar(ss []any) []any {
	for _, f := range []func(any) []any{
		desugarLetMatch,
		desugarLetExpression,
		desugarEmptyCollection,
		desugarDictionaryExpansion,
		match.Desugar,
		desugarAnonymousFunctions,
		desugarMutualRecursionStatement,
		desugarSelfRecursiveStatement,
		flattenStatement,
		removeUnusedVariables,
		removeAliases,
	} {
		new := make([]any, 0, 2*len(ss))

		for _, s := range ss {
			new = append(new, f(s)...)
		}

		ss = new
	}

	return ss
}
