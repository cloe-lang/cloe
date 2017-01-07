package vm

type Closure struct {
	function      *Thunk
	freeVariables []*Thunk
}

func (c Closure) Call(ts ...*Thunk) *Thunk {
	f, ok := c.function.Eval().(Callable)

	if !ok {
		return NewError("Something not callable was called.")
	}

	return f.Call(append(c.freeVariables, ts...)...)
}

func Partial(ts ...*Thunk) *Thunk {
	if len(ts) == 0 {
		return NewError("Number of arguments to partial must be >= 1.")
	}

	return Normal(Closure{ts[0], ts[1:]})
}
