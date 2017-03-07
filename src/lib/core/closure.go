package core

type closureType struct {
	function      *Thunk
	freeVariables Arguments
}

func (c closureType) call(args Arguments) Object {
	return App(c.function, c.freeVariables.Merge(args))
}

type RawFunction func(Arguments) Object

func (f RawFunction) call(args Arguments) Object {
	return f(args)
}

var Partial = Normal(RawFunction(func(args Arguments) Object {
	t := args.nextPositional()
	return closureType{t, args}
}))

func (c closureType) string() Object {
	return StringType("<closure>")
}
