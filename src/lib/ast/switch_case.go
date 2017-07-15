package ast

import "fmt"

// SwitchCase represents a case of a pattern and corrensponding value.
type SwitchCase struct {
	pattern string
	value   interface{}
}

// NewSwitchCase creates a case in a switch expression.
func NewSwitchCase(p string, v interface{}) SwitchCase {
	return SwitchCase{p, v}
}

// Pattern returns a pattern of a case in a switch expression.
func (c SwitchCase) Pattern() string {
	return c.pattern
}

// Value returns a value corrensponding to a pattern in a switch expression.
func (c SwitchCase) Value() interface{} {
	return c.value
}

func (c SwitchCase) String() string {
	return fmt.Sprintf("%v %v", c.pattern, c.value)
}
