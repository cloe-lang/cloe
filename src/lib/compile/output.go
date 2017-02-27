package compile

import "github.com/raviqqe/tisp/src/lib/core"

type Output struct {
	value    *core.Thunk
	expanded bool
}

func NewOutput(value *core.Thunk, expanded bool) Output {
	return Output{value, expanded}
}

func (o Output) Value() *core.Thunk {
	return o.value
}

func (o Output) Expanded() bool {
	return o.expanded
}
