package vm

type Number float64

func NewNumber(n float64) *Thunk {
	return Normal(Number(n))
}

func (n Number) Equal(e Equalable) bool {
	return n == e.(Number)
}

func (n Number) Add(a Addable) Addable {
	return n + a.(Number)
}

func Sub(ts ...*Thunk) *Thunk {
	if len(ts) == 0 {
		return NumArgsError("sub", ">= 1")
	}

	for _, t := range ts {
		go t.Eval()
	}

	o := ts[0].Eval()
	n0, ok := o.(Number)

	if !ok {
		return notNumberError(o)
	}

	for _, t := range ts[1:] {
		o := t.Eval()
		n, ok := o.(Number)

		if !ok {
			return notNumberError(o)
		}

		n0 = n0 - n
	}

	return Normal(n0)
}

func Mult(ts ...*Thunk) *Thunk {
	if len(ts) == 0 {
		// for symmetry with add while it can take no argument and return 1
		return NumArgsError("mult", ">= 1")
	}

	for _, t := range ts {
		go t.Eval()
	}

	n0 := Number(1)

	for _, t := range ts {
		o := t.Eval()
		n, ok := o.(Number)

		if !ok {
			return notNumberError(o)
		}

		n0 = n0 * n
	}

	return Normal(n0)
}

func Div(ts ...*Thunk) *Thunk {
	if len(ts) == 0 {
		return NumArgsError("div", ">= 1")
	}

	for _, t := range ts {
		go t.Eval()
	}

	o := ts[0].Eval()
	n0, ok := o.(Number)

	if !ok {
		return notNumberError(o)
	}

	for _, t := range ts[1:] {
		o := t.Eval()
		n, ok := o.(Number)

		if !ok {
			return notNumberError(o)
		}

		n0 = n0 / n
	}

	return Normal(n0)
}

func notNumberError(o Object) *Thunk {
	return TypeError(o, "Number")
}
