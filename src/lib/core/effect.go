package core

// effectType represents an effect value returned from an impure function.
// effectType is meant to be used to distinguish calls of pure and impure
// functions and also represent a "result" value of an impure function which
// can be extracted by a special function named "out" and passed to a pure
// function.
type effectType struct {
	value *Thunk
}

// newEffect creates an effect value.
func newEffect(value *Thunk) *Thunk {
	return Normal(effectType{value})
}

// Pure extracts a result value in an effect value.
var Pure = NewLazyFunction(
	NewSignature([]string{"effect"}, nil, "", nil, nil, ""),
	func(ts ...*Thunk) Value {
		return ts[0].EvalEffect()
	})

func (o effectType) string() Value {
	return StringType("<effect>")
}
