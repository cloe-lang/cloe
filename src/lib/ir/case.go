package ir

import "github.com/cloe-lang/cloe/src/lib/core"

// Case represents a case of a pattern and a corresponding value in a switch expression.
type Case struct {
	pattern core.Value
	value   any
}

// NewCase creates a case in a switch expression.
func NewCase(p core.Value, v any) Case {
	return Case{p, v}
}
