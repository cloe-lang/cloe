package ast

// Switch represents a switch expression.
type Switch struct {
	value interface{}
	cases []SwitchCase
}

// NewSwitch creates a match expression.
func NewSwitch(v interface{}, cs []SwitchCase) Switch {
	if len(cs) == 0 {
		panic("Cases in a match expression must be more than 0.")
	}

	return Switch{v, cs}
}

// Value returns a value which will be matched with patterns in a match expression.
func (s Switch) Value() interface{} {
	return s.value
}

// Cases returns pairs of a pattern and corrensponding value in a match expression.
func (s Switch) Cases() []SwitchCase {
	return s.cases
}
