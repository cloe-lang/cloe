package ast

import "fmt"

// OptionalArgument represents an optional argument defined in a function.
type OptionalArgument struct {
	name         string
	defaultValue interface{}
}

// NewOptionalArgument creates an optional argument.
func NewOptionalArgument(n string, v interface{}) OptionalArgument {
	return OptionalArgument{n, v}
}

// Name returns a name of an optional argument.
func (o OptionalArgument) Name() string {
	return o.name
}

// DefaultValue returns a default value of an optional argument.
func (o OptionalArgument) DefaultValue() interface{} {
	return o.defaultValue
}

func (o OptionalArgument) String() string {
	return fmt.Sprintf("(%v %v)", o.name, o.defaultValue)
}
