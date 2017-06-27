package desugar

import "github.com/tisp-lang/tisp/src/lib/desugar/match"

// Desugar desugars a module of statements in AST.
func Desugar(module []interface{}) []interface{} {
	new := make([]interface{}, 0)

	for _, s := range module {
		for _, s := range match.Desugar(s) {
			for _, s := range desugarMutualRecursionStatement(s) {
				for _, s := range desugarSelfRecursiveStatement(s) {
					new = append(new, flattenStatement(s)...)
				}
			}
		}
	}

	return new
}
