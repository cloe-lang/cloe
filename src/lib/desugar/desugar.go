package desugar

// Desugar desugars a statement in AST.
func Desugar(s interface{}) []interface{} {
	new := make([]interface{}, 0)

	for _, s := range desugarMutualRecursionStatement(s) {
		for _, s := range desugarSelfRecursiveStatement(s) {
			new = append(new, flattenStatement(s)...)
		}
	}

	return new
}
