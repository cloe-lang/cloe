package core

// PositionalArgument represents a positional argument.
// It can be expanded as a list.
type PositionalArgument struct {
	value    *Thunk
	expanded bool
}

// NewPositionalArgument creates a PositionalArgument.
func NewPositionalArgument(value *Thunk, expanded bool) PositionalArgument {
	return PositionalArgument{
		value:    value,
		expanded: expanded,
	}
}

// NewPositionalArguments creates an Arguments which consists of unexpanded
// positional arguments.
func NewPositionalArguments(ts ...*Thunk) Arguments {
	ps := make([]PositionalArgument, 0, len(ts))

	for _, t := range ts {
		ps = append(ps, NewPositionalArgument(t, false))
	}

	return NewArguments(ps, nil, nil)
}
