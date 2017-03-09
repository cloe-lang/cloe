package compile

// Compile compiles AST into outputs of thunks.
func Compile(module []interface{}) []Output {
	c := newCompiler()
	return c.compile(module)
}
