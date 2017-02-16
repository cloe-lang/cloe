package sig

import ".."

type PositionalArgument struct {
	value    *vm.Thunk
	expanded bool
}

func NewPositionalArgument(value *vm.Thunk, expanded bool) PositionalArgument {
	return PositionalArgument{
		value:    value,
		expanded: expanded,
	}
}
