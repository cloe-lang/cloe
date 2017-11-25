package core

type rawFunction func(Arguments) Value

func (f rawFunction) call(args Arguments) Value {
	return f(args)
}

// Partial creates a partially-applied function with arguments.
var Partial = Normal(rawFunction(func(vars Arguments) Value {
	return Normal(rawFunction(func(args Arguments) Value {
		vars := vars
		t := vars.nextPositional()
		return App(t, vars.Merge(args))
	}))
}))

func (f rawFunction) string() Value {
	return StringType("<function>")
}
