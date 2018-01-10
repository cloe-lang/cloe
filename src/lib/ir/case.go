package ir

import "github.com/coel-lang/coel/src/lib/core"

// Case represents a case of a pattern and a corresponding value in a switch expression.
type Case struct {
	pattern core.Value
	value   interface{}
}

// NewCase creates a case in a switch expression.
func NewCase(p core.Value, v interface{}) Case {
	return Case{p, v}
}
