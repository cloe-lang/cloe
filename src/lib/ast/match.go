package ast

// Match represents a match expression.
type Match struct {
	value interface{}
	cases []Case
}

// NewMatch creates a match expression.
func NewMatch(v interface{}, cs []Case) Match {
	if len(cs) == 0 {
		panic("Cases in a match expression must be more than 0.")
	}

	return Match{v, cs}
}

// Value returns a value which will be matched with patterns in a match expression.
func (m Match) Value() interface{} {
	return m.value
}

// Cases returns pairs of a pattern and corrensponding value in a match expression.
func (m Match) Cases() []Case {
	return m.cases
}
