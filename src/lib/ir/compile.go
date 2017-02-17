package ir

import (
	"../vm"
	"fmt"
)

func Compile(s vm.Signature, expr interface{}) *vm.Thunk {
	return vm.NewLazyFunction(
		s,
		func(ts ...*vm.Thunk) vm.Object {
			return compileExpression(ts, expr)
		})
}

func compileExpression(args []*vm.Thunk, expr interface{}) *vm.Thunk {
	switch x := expr.(type) {
	case int:
		return args[x]
	case *vm.Thunk:
		return x
	case IRThunk:
		return x.compile(args)
	}

	panic(fmt.Sprintf("Invalid type. %v", expr))
}
