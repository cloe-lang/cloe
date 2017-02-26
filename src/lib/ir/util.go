package ir

import (
	"../core"
	"fmt"
)

func CompileFunction(s core.Signature, expr interface{}) *core.Thunk {
	return core.NewLazyFunction(
		s,
		func(ts ...*core.Thunk) core.Object {
			return compileExpression(ts, expr)
		})
}

func compileExpression(args []*core.Thunk, expr interface{}) *core.Thunk {
	switch x := expr.(type) {
	case int:
		return args[x]
	case *core.Thunk:
		return x
	case App:
		return x.compile(args)
	}

	panic(fmt.Sprintf("Invalid type. %v", expr))
}
