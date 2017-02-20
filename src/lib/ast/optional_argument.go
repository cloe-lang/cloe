package ast

type OptionalArgument struct {
	name         string
	defaultValue interface{}
}

func NewOptionalArgument(n string, v interface{}) OptionalArgument {
	return OptionalArgument{n, v}
}

func (o OptionalArgument) Name() string {
	return o.name
}

func (o OptionalArgument) DefaultValue() interface{} {
	return o.defaultValue
}
