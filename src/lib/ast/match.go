package ast

import (
	"fmt"
	"strings"
)

// Match represents a match expression.
type Match struct {
	value interface{}
	cases []MatchCase
}

// NewMatch creates a match expression.
func NewMatch(v interface{}, cs []MatchCase) Match {
	if len(cs) == 0 {
		panic("MatchCases in a match expression must be more than 0.")
	}

	return Match{v, cs}
}

// Value returns a value which will be matched with patterns in a match expression.
func (m Match) Value() interface{} {
	return m.value
}

// Cases returns pairs of a pattern and corrensponding value in a match expression.
func (m Match) Cases() []MatchCase {
	return m.cases
}

func (m Match) String() string {
	ss := make([]string, 0, len(m.cases))

	for _, c := range m.cases {
		ss = append(ss, c.String())
	}

	return fmt.Sprintf("(match %v %v)", m.value, strings.Join(ss, " "))
}
