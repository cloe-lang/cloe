package ast

// Case represents a case of a pattern and corrensponding value.
type Case struct {
	pattern interface{}
	value   interface{}
}

// NewCase creates a case in a match expression.
func NewCase(p interface{}, v interface{}) Case {
	return Case{p, v}
}

// Pattern returns a pattern of a case in a match expression.
func (c Case) Pattern() interface{} {
	return c.pattern
}

// Value returns a value corrensponding to a pattern in a match expression.
func (c Case) Value() interface{} {
	return c.value
}
