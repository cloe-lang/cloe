package compile

import "github.com/coel-lang/coel/src/lib/core"

// Effect represents an effect of a program.
type Effect struct {
	value    *core.Thunk
	expanded bool
}

// NewEffect creates an effect.
func NewEffect(value *core.Thunk, expanded bool) Effect {
	return Effect{value, expanded}
}

// Value returns an effect of a thunk.
func (o Effect) Value() *core.Thunk {
	return o.value
}

// Expanded returns true if it is a expanded list of effects or false otherwise.
func (o Effect) Expanded() bool {
	return o.expanded
}
