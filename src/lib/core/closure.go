package core

type closureType struct {
	function      *Thunk
	freeVariables Arguments
}

func (c closureType) call(args Arguments) Object {
	return App(c.function, c.freeVariables.Merge(args))
}

type rawFunction func(Arguments) Object

func (f rawFunction) call(args Arguments) Object {
	return f(args)
}

// Partial creates a partially-applied function with arguments.
var Partial = Normal(rawFunction(func(args Arguments) Object {
	t := args.nextPositional()
	return closureType{t, args}
}))

func (c closureType) string() Object {
	return StringType("<closure>")
}
