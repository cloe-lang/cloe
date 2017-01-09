package vm

type Closure struct {
	function      *Thunk
	freeVariables []*Thunk
}

func (c Closure) Call(ts ...*Thunk) Object {
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

	return Closure{ts[0], ts[1:]}
}
