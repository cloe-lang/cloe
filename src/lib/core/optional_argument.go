package core

// OptionalArgument represents an optional argument defined in a function.
type OptionalArgument struct {
	name         string
	defaultValue Value
}

// NewOptionalArgument creates an optional argument.
func NewOptionalArgument(n string, v Value) OptionalArgument {
	return OptionalArgument{n, v}
}
