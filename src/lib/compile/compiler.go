package compile

import (
	"../ast"
	"../ir"
	"../vm"
	"./env"
	"fmt"
)

type compiler struct {
	env     env.Environment
	outputs []*vm.Thunk
}

func newCompiler() compiler {
	return compiler{
		env:     prelude.Child(),
		outputs: make([]*vm.Thunk, 0, 8), // TODO: Best cap?
	}
}

func (c *compiler) compile(module []interface{}) []*vm.Thunk {
	for _, node := range module {
		switch x := node.(type) {
		case ast.LetConst:
			c.env.Set(x.Name(), c.compileExpression(x.Expr()))
		case ast.LetFunction:
			c.env.Set(x.Name(), ir.CompileFunction(x.Signature(), c.compileFunctionBodyToIR(x.Body())))
		case ast.Output:
			c.outputs = append(c.outputs, c.compileExpression(x.Expr()))
		default:
			panic(fmt.Sprint("Invalid instruction.", x))
		}
	}

	return c.outputs
}

func (c *compiler) compileExpression(expr interface{}) *vm.Thunk {
	switch x := expr.(type) {
	case string:
		return getOrError(c.env, x)
	case []interface{}:
		ts := make([]*vm.Thunk, len(x))

		for i, e := range x {
			ts[i] = c.compileExpression(e)
		}

		return vm.PApp(ts[0], ts[1:]...)
	}

	panic(fmt.Sprint("Invalid type as an expression.", expr))
}

func (c *compiler) compileFunctionBodyToIR(expr interface{}) interface{} {
	switch x := expr.(type) {
	case string:
		return getOrError(c.env, x)
	case int:
		return x
	case []interface{}:
		ps := make([]ir.PositionalArgument, len(x)-1)

		for i, e := range x[1:] {
			ps[i] = ir.NewPositionalArgument(c.compileFunctionBodyToIR(e), false)
		}

		// TODO: Support keyword arguments and expanded dictionaries.
		return ir.NewApp(
			c.compileFunctionBodyToIR(x[0]),
			ir.NewArguments(ps, []ir.KeywordArgument{}, []interface{}{}))
	}

	panic(fmt.Sprint("Invalid type.", expr))
}

func getOrError(e env.Environment, s string) *vm.Thunk {
	t, err := e.Get(s)

	if err != nil {
		panic(err)
	}

	return t
}
