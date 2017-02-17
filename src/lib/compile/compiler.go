package compile

import (
	"../vm"
	"fmt"
)

type compiler struct {
	env     *environment
	outputs []*vm.Thunk
}

func newCompiler() *compiler {
	return &compiler{
		env:     newEnvironment(nil),
		outputs: make([]*vm.Thunk, 0, 8), // TODO: Best cap?
	}
}

func (c *compiler) compile(instrs []interface{}) []*vm.Thunk {
	for _, instr := range instrs {
		switch x := instr.(type) {
		case LetConst:
			c.env.set(x.name, c.compileExpression(x.expr))
		case LetFunction:
			c.env.set(x.name, compileFunction(x.signature, c.compileFunctionBodyToIR(x.body)))
		case Output:
			c.outputs = append(c.outputs, c.compileExpression(x.expr))
		default:
			panic(fmt.Sprint("Invalid instruction.", x))
		}
	}

	return c.outputs
}

func (c *compiler) compileExpression(e Expression) *vm.Thunk {
	switch x := e.(type) {
	case string:
		return c.env.get(x)
	case []interface{}:
		ts := make([]*vm.Thunk, len(x))

		for i, e := range x {
			ts[i] = c.compileExpression(e)
		}

		return vm.PApp(ts[0], ts[1:]...)
	}

	panic(fmt.Sprint("Invalid type as an expression.", e))
}

func (c *compiler) compileFunctionBodyToIR(e Expression) interface{} {
	switch x := e.(type) {
	case string:
		return c.env.get(x)
	case int:
		return x
	case []interface{}:
		ps := make([]PositionalArgument, len(x)-1)

		for i, e := range x[1:] {
			ps[i] = NewPositionalArgument(c.compileFunctionBodyToIR(e), false)
		}

		// TODO: Support keyword arguments and expanded dictionaries.
		return App(
			c.compileFunctionBodyToIR(x[0]),
			NewArguments(ps, []KeywordArgument{}, []interface{}{}))
	}

	panic(fmt.Sprint("Invalid type.", e))
}

func compileFunction(s vm.Signature, expr interface{}) *vm.Thunk {
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
	case Thunk:
		return x.compile(args)
	}

	panic(fmt.Sprintf("Invalid type. %v", expr))
}
