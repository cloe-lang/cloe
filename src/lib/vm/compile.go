package vm

import "fmt"

func Compile(s Signature, expr interface{}) *Thunk {
	return NewLazyFunction(
		s,
		func(ts ...*Thunk) Object {
			return compileExpression(ts, expr)
		})
}

func compileExpression(args []*Thunk, expr interface{}) *Thunk {
	switch x := expr.(type) {
	case int:
		return args[x]
	case *Thunk:
		return x
	case IRThunk:
		return x.compile(args)
	}

	panic(fmt.Sprintf("Invalid type. %v", expr))
}
