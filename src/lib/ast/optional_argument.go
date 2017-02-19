package ast

type OptionalArgument struct {
	name  string
	value interface{}
}

func NewOptionalArgument(n string, v interface{}) OptionalArgument {
	return OptionalArgument{n, v}
}
