package core

// effectType represents an effect value returned from an impure function.
// effectType is meant to be used to distinguish calls of pure and impure
// functions and also represent a "result" value of an impure function which
// can be extracted by a special function named "out" and passed to a pure
// function.
type effectType struct {
	value Value
}

// Eval evaluates a value into a WHNF.
func (e effectType) eval() Value {
	return e
}

// newEffect creates an effect value.
func newEffect(value Value) Value {
	return effectType{value}
}

// Pure extracts a result value in an effect value.
var Pure = NewLazyFunction(
	NewSignature([]string{"arg"}, nil, "", nil, nil, ""),
	func(vs ...Value) Value {
		return EvalImpure(vs[0])
	})
