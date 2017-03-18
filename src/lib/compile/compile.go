package compile

import "github.com/raviqqe/tisp/src/lib/core"

// MainModule compiles an AST into outputs of thunks.
func MainModule(module []interface{}) []Output {
	c := newCompiler()
	return c.compile(module)
}

// SubModule compiles an AST into a map of names to thunks..
func SubModule(module []interface{}) map[string]*core.Thunk {
	c := newCompiler()
	c.compile(module)
	return c.env.toMap()
}
