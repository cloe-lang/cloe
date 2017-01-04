package functions

import "../vm"

func Nil() *vm.Thunk {
	return vm.Normal(nil)
}
