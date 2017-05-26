package core

// OptionalArgument represents an optional argument defined in a function.
type OptionalArgument struct {
	name         string
	defaultValue *Thunk
}

// NewOptionalArgument creates an optional argument.
func NewOptionalArgument(n string, v *Thunk) OptionalArgument {
	return OptionalArgument{
		name:         n,
		defaultValue: v,
	}
}
