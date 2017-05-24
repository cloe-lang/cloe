package compile

import "github.com/tisp-lang/tisp/src/lib/core"

// Output represents an output of a program.
type Output struct {
	value    *core.Thunk
	expanded bool
}

// NewOutput creates an output.
func NewOutput(value *core.Thunk, expanded bool) Output {
	return Output{value, expanded}
}

// Value returns an output of a thunk.
func (o Output) Value() *core.Thunk {
	return o.value
}

// Expanded returns true if it is a expanded list of outputs or false otherwise.
func (o Output) Expanded() bool {
	return o.expanded
}
