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
