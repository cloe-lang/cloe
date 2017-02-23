package ast

type PositionalArgument struct {
	value    interface{}
	expanded bool
}

func NewPositionalArgument(value interface{}, expanded bool) PositionalArgument {
	return PositionalArgument{value, expanded}
}

func (p PositionalArgument) Value() interface{} {
	return p.value
}

func (p PositionalArgument) Expanded() bool {
	return p.expanded
}
