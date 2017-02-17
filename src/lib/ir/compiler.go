package ir

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
			c.env.set(x.name, vm.Compile(x.signature, c.compileFunctionBodyToIR(x.body)))
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
		ps := make([]vm.IRPositionalArgument, len(x)-1)

		for i, e := range x[1:] {
			ps[i] = vm.NewIRPositionalArgument(c.compileFunctionBodyToIR(e), false)
		}

		// TODO: Support keyword arguments and expanded dictionaries.
		return vm.IRApp(
			c.compileFunctionBodyToIR(x[0]),
			vm.NewIRArguments(ps, []vm.IRKeywordArgument{}, []interface{}{}))
	}

	panic(fmt.Sprint("Invalid type.", e))
}
