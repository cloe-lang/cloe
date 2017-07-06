package ast

// MatchCase represents a case of a pattern and corrensponding value.
type MatchCase struct {
	pattern interface{}
	value   interface{}
}

// NewMatchCase creates a case in a match expression.
func NewMatchCase(p interface{}, v interface{}) MatchCase {
	return MatchCase{p, v}
}

// Pattern returns a pattern of a case in a match expression.
func (c MatchCase) Pattern() interface{} {
	return c.pattern
}

// Value returns a value corrensponding to a pattern in a match expression.
func (c MatchCase) Value() interface{} {
	return c.value
}
