package ir

import (
	"../vm"
	"fmt"
)

var write = vm.NewStrictFunction(
	vm.NewSignature(
		[]string{"x"}, []vm.OptionalArgument{}, "",
		[]string{}, []vm.OptionalArgument{}, "",
	),
	func(os ...vm.Object) vm.Object {
		fmt.Println(os[0])

		return vm.Nil
	})
