package core

// OutputType represents an output value returned from an impure function.
type OutputType struct {
	value *Thunk
}

// NewOutput creates an output value.
func NewOutput(value *Thunk) *Thunk {
	return Normal(OutputType{value})
}
