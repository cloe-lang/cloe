package ir

import (
	"../vm"
	"fmt"
)

var write = vm.NewStrictFunction(func(os ...vm.Object) vm.Object {
	if len(os) != 1 {
		return vm.NumArgsError("write", "1")
	}

	fmt.Println(os[0])

	return vm.Nil
})
