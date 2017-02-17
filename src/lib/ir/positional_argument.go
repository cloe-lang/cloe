package ir

import "../vm"

type PositionalArgument struct {
	value    interface{}
	expanded bool
}

func NewPositionalArgument(ir interface{}, expanded bool) PositionalArgument {
	return PositionalArgument{
		value:    ir,
		expanded: expanded,
	}
}

func (p PositionalArgument) compile(args []*vm.Thunk) vm.PositionalArgument {
	return vm.NewPositionalArgument(compileExpression(args, p.value), p.expanded)
}
