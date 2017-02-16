package sig

import ".."

// OptionalArgument contains its name and default value.
type OptionalArgument struct {
	name         string
	defaultValue *vm.Thunk
}

// NewOptionalArgument defines a new OptionalArgument.
func NewOptionalArgument(n string, v *vm.Thunk) OptionalArgument {
	return OptionalArgument{
		name:         n,
		defaultValue: v,
	}
}
