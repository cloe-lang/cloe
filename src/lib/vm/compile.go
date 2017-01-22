package vm

import "fmt"

// alias expr = int | *Thunk | []expr
func Compile(expr interface{}) *Thunk {
	return NewLazyFunction(func(ts ...*Thunk) Object {
		return compileExpression(ts, expr)
	})
}

func compileExpression(args []*Thunk, expr interface{}) *Thunk {
	switch x := expr.(type) {
	case int:
		return args[x]
	case *Thunk:
		return x
	case []interface{}:
		ts := make([]*Thunk, 0, len(x))

		for _, expr := range x {
			ts = append(ts, compileExpression(args, expr))
		}

		return App(ts[0], ts[1:]...)
	}

	panic(fmt.Sprintf("Invalid type. {}", expr))
}
