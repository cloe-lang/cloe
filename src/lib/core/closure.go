package core

type RawFunction func(Arguments) Value

func (f RawFunction) call(args Arguments) Value {
	return f(args)
}

// Partial creates a partially-applied function with arguments.
var Partial = Normal(RawFunction(func(vars Arguments) Value {
	return Normal(RawFunction(func(args Arguments) Value {
		vars := vars
		t := vars.nextPositional()
		return App(t, vars.Merge(args))
	}))
}))

func (f RawFunction) string() Value {
	return StringType("<function>")
}
