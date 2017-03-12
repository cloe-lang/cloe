package core

type closureType struct {
	function      *Thunk
	freeVariables Arguments
}

func (c closureType) call(args Arguments) Value {
	return App(c.function, c.freeVariables.Merge(args))
}

type rawFunction func(Arguments) Value

func (f rawFunction) call(args Arguments) Value {
	return f(args)
}

// Partial creates a partially-applied function with arguments.
var Partial = Normal(rawFunction(func(args Arguments) Value {
	t := args.nextPositional()
	return closureType{t, args}
}))

func (c closureType) string() Value {
	return StringType("<closure>")
}
