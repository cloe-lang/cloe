package core

// EffectType represents an effect value returned from an impure function.
// EffectType is meant to be used to distinguish calls of pure and impure
// functions and also represent a "result" value of an impure function which
// can be extracted by a special function named "out" and passed to a pure
// function.
type EffectType struct {
	value *Thunk
}

// NewEffect creates an effect value.
func NewEffect(value *Thunk) *Thunk {
	return Normal(EffectType{value})
}

// Pure extracts a result value in an effect value.
var Pure = NewLazyFunction(
	NewSignature([]string{"effect"}, nil, "", nil, nil, ""),
	func(ts ...*Thunk) Value {
		return ts[0].EvalEffect()
	})

func (o EffectType) string() Value {
	return StringType("<effect>")
}
