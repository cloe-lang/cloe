package ast

type PositionalArgument struct {
	value    interface{}
	expanded bool
}

func NewPositionalArgument(value interface{}, expanded bool) PositionalArgument {
	return PositionalArgument{value, expanded}
}
