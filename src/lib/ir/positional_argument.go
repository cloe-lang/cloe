package ir

import "github.com/tisp-lang/tisp/src/lib/core"

// PositionalArgument represents a positional argument passed to a function.
// It can be a list value and expanded into multiple arguments.
type PositionalArgument struct {
	value    interface{}
	expanded bool
}

// NewPositionalArgument creates a positional argument.
func NewPositionalArgument(ir interface{}, expanded bool) PositionalArgument {
	return PositionalArgument{
		value:    ir,
		expanded: expanded,
	}
}

func (p PositionalArgument) interpret(args []*core.Thunk) core.PositionalArgument {
	return core.NewPositionalArgument(interpretExpression(args, p.value), p.expanded)
}
