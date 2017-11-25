package core

type RawFunctionType func(Arguments) Value

func (f RawFunctionType) call(args Arguments) Value {
	return f(args)
}

// Partial creates a partially-applied function with arguments.
var Partial = Normal(RawFunctionType(func(vars Arguments) Value {
	return Normal(RawFunctionType(func(args Arguments) Value {
		vars := vars
		t := vars.nextPositional()
		return App(t, vars.Merge(args))
	}))
}))

func (f RawFunctionType) string() Value {
	return StringType("<function>")
}
