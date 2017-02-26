package core

// OptionalArgument contains its name and default value.
type OptionalArgument struct {
	name         string
	defaultValue *Thunk
}

// NewOptionalArgument defines a new OptionalArgument.
func NewOptionalArgument(n string, v *Thunk) OptionalArgument {
	return OptionalArgument{
		name:         n,
		defaultValue: v,
	}
}
