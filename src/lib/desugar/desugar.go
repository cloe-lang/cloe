package desugar

// Desugar desugars a module of statements in AST.
func Desugar(module []interface{}) []interface{} {
	new := make([]interface{}, 0, 2*len(module)) // TODO: Best cap?

	for _, s := range module {
		new = append(new, flattenStatement(s)...)
	}

	return new
}
