package vm

type closureType struct {
	function      *Thunk
	freeVariables []*Thunk
}

func (c closureType) call(ts ...*Thunk) Object {
	o := c.function.Eval()
	f, ok := o.(callable)

	if !ok {
		return notCallableError(o)
	}

	return f.call(append(c.freeVariables, ts...)...)
}

var Partial = NewLazyFunction(partial)

func partial(ts ...*Thunk) Object {
	switch len(ts) {
	case 0:
		return numArgsError("partial", ">= 1")
	case 1:
		return ts[0]
	}

	return closureType{ts[0], ts[1:]}
}
