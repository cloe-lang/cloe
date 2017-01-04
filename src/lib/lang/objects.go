package functions

import "../vm"

func Nil() *vm.Thunk {
	return vm.Normal(nil)
}

func Number(n float64) *vm.Thunk {
	return vm.Normal(vm.NewNumber(n))
}
