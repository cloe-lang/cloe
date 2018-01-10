package core

// PositionalArgument represents a positional argument.
// It can be expanded as a list.
type PositionalArgument struct {
	value    Value
	expanded bool
}

// NewPositionalArgument creates a PositionalArgument.
func NewPositionalArgument(v Value, expanded bool) PositionalArgument {
	return PositionalArgument{v, expanded}
}
