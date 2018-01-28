package ir

import "github.com/coel-lang/coel/src/lib/core"

// CompileFunction compiles a function in IR into a thunk.
func CompileFunction(s core.Signature, vars []interface{}, expr interface{}) core.Value {
	// TODO: Compile everything into bytecode here.

	return core.NewLazyFunction(
		s,
		func(ts ...core.Value) core.Value {
			args := append(make([]core.Value, 0, len(ts)+len(vars)), ts...)

			for _, v := range vars {
				args = append(args, interpretExpression(args, v))
			}

			return interpretExpression(args, expr)
		})
}

func interpretExpression(args []core.Value, expr interface{}) core.Value {
	switch x := expr.(type) {
	case int:
		return args[x]
	case core.Value:
		return x
	case App:
		return x.interpret(args)
	case Switch:
		v := core.EvalPure(core.PApp(x.dict, interpretExpression(args, x.matchedValue)))
		n, ok := v.(*core.NumberType)

		if !ok {
			return interpretExpression(args, x.defaultCase)
		}

		return interpretExpression(args, x.caseValues[int(*n)])
	}

	panic("Unreachable")
}
