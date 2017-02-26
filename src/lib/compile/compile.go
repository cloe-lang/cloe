package compile

import "../core"

func Compile(module []interface{}) []*core.Thunk {
	c := newCompiler()
	return c.compile(module)
}
