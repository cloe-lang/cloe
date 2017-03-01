package ir

import (
	"fmt"
	"github.com/raviqqe/tisp/src/lib/core"
)

func CompileFunction(s core.Signature, vars []interface{}, expr interface{}) *core.Thunk {
	return core.NewLazyFunction(
		s,
		func(ts ...*core.Thunk) core.Object {
			return compileWithVars(ts, vars, expr)
		})
}

func compileWithVars(args []*core.Thunk, vars []interface{}, expr interface{}) *core.Thunk {
	if len(vars) == 0 {
		return compileExpression(args, expr)
	}

	return compileWithVars(
		append(args, compileExpression(args, vars[0])),
		vars[1:],
		expr)
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
