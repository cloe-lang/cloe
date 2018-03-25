package core

// OptionalParameter represents an optional argument defined in a function.
type OptionalParameter struct {
	name         string
	defaultValue Value
}

// NewOptionalParameter creates an optional argument.
func NewOptionalParameter(n string, v Value) OptionalParameter {
	return OptionalParameter{n, v}
}
