package vm

type closureType struct {
	function      *Thunk
	freeVariables []*Thunk
}

func (c closureType) Call(ts ...*Thunk) Object {
	o := c.function.Eval()
	f, ok := o.(Callable)

	if !ok {
		return NotCallableError(o)
	}

	return f.Call(append(c.freeVariables, ts...)...)
}

var Partial = NewLazyFunction(partial)

func partial(ts ...*Thunk) Object {
	switch len(ts) {
	case 0:
		return NumArgsError("partial", ">= 1")
	case 1:
		return ts[0]
	}

	return closureType{ts[0], ts[1:]}
}
