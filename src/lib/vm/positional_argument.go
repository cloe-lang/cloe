package vm

type PositionalArgument struct {
	value    *Thunk
	expanded bool
}

func NewPositionalArgument(value *Thunk, expanded bool) PositionalArgument {
	return PositionalArgument{
		value:    value,
		expanded: expanded,
	}
}
