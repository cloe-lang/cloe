package vm

type Number float64

func NewNumber(n float64) Number {
	return Number(n)
}

func NumberThunk(n float64) *Thunk {
	return Normal(NewNumber(n))
}

func (n Number) Equal(e Equalable) Object {
	return NewBool(n == e.(Number))
}

func (n Number) Add(a Addable) Addable {
	return n + a.(Number)
}

var Sub = NewStrictFunction(sub)

func sub(os ...Object) Object {
	if len(os) == 0 {
		return NumArgsError("sub", ">= 1")
	}

	o := os[0]
	n0, ok := o.(Number)

	if !ok {
		return notNumberError(o)
	}

	for _, o := range os[1:] {
		n, ok := o.(Number)

		if !ok {
			return notNumberError(o)
		}

		n0 -= n
	}

	return n0
}

var Mult = NewStrictFunction(mult)

func mult(os ...Object) Object {
	if len(os) == 0 {
		// for symmetry with add while it can take no argument and return 1.
		return NumArgsError("mult", ">= 1")
	}

	n0 := Number(1)

	for _, o := range os {
		n, ok := o.(Number)

		if !ok {
			return notNumberError(o)
		}

		n0 *= n
	}

	return n0
}

var Div = NewStrictFunction(div)

func div(os ...Object) Object {
	if len(os) == 0 {
		return NumArgsError("div", ">= 1")
	}

	o := os[0]
	n0, ok := o.(Number)

	if !ok {
		return notNumberError(o)
	}

	for _, o := range os[1:] {
		n, ok := o.(Number)

		if !ok {
			return notNumberError(o)
		}

		n0 /= n
	}

	return n0
}

func notNumberError(o Object) Error {
	return TypeError(o, "Number")
}
