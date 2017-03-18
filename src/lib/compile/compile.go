package compile

// MainModule compiles AST into outputs of thunks.
func MainModule(module []interface{}) []Output {
	c := newCompiler()
	return c.compile(module)
}
