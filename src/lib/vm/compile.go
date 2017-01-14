package vm

import "fmt"

func Compile(ir interface{}) *Thunk {
	return NewLazyFunction(func(ts ...*Thunk) Object {
		return compileIR(ts, ir)
	})
}

func compileIR(args []*Thunk, ir interface{}) *Thunk {
	switch x := ir.(type) {
	case int:
		return args[x]
	case []interface{}:
		ts := make([]*Thunk, 0, len(x))

		for _, expr := range x {
			ts = append(ts, compileIR(args, expr))
		}

		return App(ts[0], ts[1:]...)
	default:
		panic(fmt.Sprintf("Invalid type is found in IR. {}", x))
	}
}
