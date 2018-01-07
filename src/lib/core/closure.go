package core

// RawFunctionType represents a raw function which takes an arguments struct
// and returns a value.
type RawFunctionType func(Arguments) Value

func (f RawFunctionType) call(args Arguments) Value {
	return f(args)
}

// Partial creates a partially-applied function with arguments.
var Partial = Normal(RawFunctionType(func(vars Arguments) Value {
	return Normal(RawFunctionType(func(args Arguments) Value {
		vars := vars
		v := vars.nextPositional().Eval()
		f, ok := v.(callable)

		if !ok {
			return TypeError(v, "callable")
		}

		return f.call(vars.Merge(args))
	}))
}))

func (f RawFunctionType) string() Value {
	return StringType("<function>")
}
