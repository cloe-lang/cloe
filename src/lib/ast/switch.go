package ast

import (
	"fmt"
	"strings"
)

// Switch represents a switch expression.
type Switch struct {
	value       any
	cases       []SwitchCase
	defaultCase any
}

// NewSwitch creates a match expression.
func NewSwitch(v any, cs []SwitchCase, d any) Switch {
	if len(cs) == 0 && d == nil {
		panic("Cases in a match expression must be more than 0.")
	}

	return Switch{v, cs, d}
}

// Value returns a value which will be matched with patterns in a match expression.
func (s Switch) Value() any {
	return s.value
}

// Cases returns pairs of a pattern and corresponding value in a match expression.
func (s Switch) Cases() []SwitchCase {
	return s.cases
}

// DefaultCase returns a default case in a switch expression.
func (s Switch) DefaultCase() any {
	return s.defaultCase
}

func (s Switch) String() string {
	ss := make([]string, 0, len(s.cases))

	for _, c := range s.cases {
		ss = append(ss, c.String())
	}

	return fmt.Sprintf("(switch %v %v %v)", s.value, strings.Join(ss, " "), s.defaultCase)
}
