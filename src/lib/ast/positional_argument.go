package ast

import "fmt"

// PositionalArgument represents a positional argument passed to a function.
type PositionalArgument struct {
	value    interface{}
	expanded bool
}

// NewPositionalArgument creates a positional argument.
func NewPositionalArgument(value interface{}, expanded bool) PositionalArgument {
	return PositionalArgument{value, expanded}
}

// Value returns a value of a positional argument.
func (p PositionalArgument) Value() interface{} {
	return p.value
}

// Expanded returns true if a positional argument is an expanded list or false otherwise.
func (p PositionalArgument) Expanded() bool {
	return p.expanded
}

func (p PositionalArgument) String() string {
	if p.expanded {
		return fmt.Sprintf("..%v", p.value)
	}

	return fmt.Sprintf("%v", p.value)
}
