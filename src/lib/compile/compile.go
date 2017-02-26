package compile

import "github.com/raviqqe/tisp/src/lib/core"

func Compile(module []interface{}) []*core.Thunk {
	c := newCompiler()
	return c.compile(module)
}
