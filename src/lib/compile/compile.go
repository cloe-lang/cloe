package compile

import "../vm"

func Compile(module []interface{}) []*vm.Thunk {
	c := newCompiler()
	return c.compile(module)
}
