package vm

type numberType float64

func NewNumber(n float64) *Thunk {
	return Normal(numberType(n))
}

func (n numberType) equal(e equalable) Object {
	return rawBool(n == e.(numberType))
}

func (n numberType) add(a addable) addable {
	return n + a.(numberType)
}

var Sub = NewStrictFunction(sub)

func sub(os ...Object) Object {
	if len(os) == 0 {
		return numArgsError("sub", ">= 1")
	}

	o := os[0]
	n0, ok := o.(numberType)

	if !ok {
		return notNumberError(o)
	}

	for _, o := range os[1:] {
		n, ok := o.(numberType)

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
		return numArgsError("mult", ">= 1")
	}

	n0 := numberType(1)

	for _, o := range os {
		n, ok := o.(numberType)

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
		return numArgsError("div", ">= 1")
	}

	o := os[0]
	n0, ok := o.(numberType)

	if !ok {
		return notNumberError(o)
	}

	for _, o := range os[1:] {
		n, ok := o.(numberType)

		if !ok {
			return notNumberError(o)
		}

		n0 /= n
	}

	return n0
}

func notNumberError(o Object) *Thunk {
	return typeError(o, "Number")
}
