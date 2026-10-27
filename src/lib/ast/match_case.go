package ast

import "fmt"

// MatchCase represents a case of a pattern and corrensponding value.
type MatchCase struct {
	pattern any
	value   any
}

// NewMatchCase creates a case in a match expression.
func NewMatchCase(p any, v any) MatchCase {
	return MatchCase{p, v}
}

// Pattern returns a pattern of a case in a match expression.
func (c MatchCase) Pattern() any {
	return c.pattern
}

// Value returns a value corrensponding to a pattern in a match expression.
func (c MatchCase) Value() any {
	return c.value
}

func (c MatchCase) String() string {
	return fmt.Sprintf("%v %v", c.pattern, c.value)
}
