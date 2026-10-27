package ast

import "fmt"

// SwitchCase represents a case of a pattern and corresponding value.
type SwitchCase struct {
	pattern string
	value   any
}

// NewSwitchCase creates a case in a switch expression.
func NewSwitchCase(p string, v any) SwitchCase {
	return SwitchCase{p, v}
}

// Pattern returns a pattern of a case in a switch expression.
func (c SwitchCase) Pattern() string {
	return c.pattern
}

// Value returns a value corresponding to a pattern in a switch expression.
func (c SwitchCase) Value() any {
	return c.value
}

func (c SwitchCase) String() string {
	return fmt.Sprintf("%v %v", c.pattern, c.value)
}
