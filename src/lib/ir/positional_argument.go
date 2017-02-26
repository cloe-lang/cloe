package ir

import "../core"

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

func (p PositionalArgument) compile(args []*core.Thunk) core.PositionalArgument {
	return core.NewPositionalArgument(compileExpression(args, p.value), p.expanded)
}
