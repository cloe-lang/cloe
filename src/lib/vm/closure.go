package vm

type Closure struct {
	function      *Thunk
	freeVariables []*Thunk
}

func (c Closure) Call(ts ...*Thunk) *Thunk {
	o := c.function.Eval()
	f, ok := o.(Callable)

	if !ok {
		return NotCallableError(o)
	}

	return f.Call(append(c.freeVariables, ts...)...)
}

func Partial(ts ...*Thunk) *Thunk {
	switch len(ts) {
	case 0:
		return NumArgsError("partial", ">= 1")
	case 1:
		return ts[0]
	}

	return Normal(Closure{ts[0], ts[1:]})
}
