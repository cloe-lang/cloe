package ast

import "fmt"

// OptionalParameter represents an optional argument defined in a function.
type OptionalParameter struct {
	name         string
	defaultValue interface{}
}

// NewOptionalParameter creates an optional argument.
func NewOptionalParameter(n string, v interface{}) OptionalParameter {
	return OptionalParameter{n, v}
}

// Name returns a name of an optional argument.
func (o OptionalParameter) Name() string {
	return o.name
}

// DefaultValue returns a default value of an optional argument.
func (o OptionalParameter) DefaultValue() interface{} {
	return o.defaultValue
}

func (o OptionalParameter) String() string {
	return fmt.Sprintf("(%v %v)", o.name, o.defaultValue)
}
