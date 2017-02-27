package compile

func Compile(module []interface{}) []Output {
	c := newCompiler()
	return c.compile(module)
}
