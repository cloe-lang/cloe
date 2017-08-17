package core

// OutputType represents an output value returned from an impure function.
// OutputType is meant to be used to distinguish calls of pure and impure
// functions and also represent a "result" value of an impure function which
// can be extracted by a special function named "out" and passed to a pure
// function.
type OutputType struct {
	value *Thunk
}

// NewOutput creates an output value.
func NewOutput(value *Thunk) *Thunk {
	return Normal(OutputType{value})
}

// Pure extracts a result value in an output value.
var Pure = NewLazyFunction(
	NewSignature([]string{"output"}, nil, "", nil, nil, ""),
	func(ts ...*Thunk) Value {
		return ts[0].EvalOutput()
	})
