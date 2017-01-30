package vm

import "math"

type numberType float64

func NewNumber(n float64) *Thunk {
	return Normal(numberType(n))
}

func (n numberType) equal(e equalable) Object {
	return rawBool(n == e.(numberType))
}

var Add = NewStrictFunction(func(os ...Object) Object {
	sum := numberType(0)

	for _, o := range os {
		n, ok := o.(numberType)

		if !ok {
			return notNumberError(o)
		}

		sum += n
	}

	return sum
})

var Sub = NewStrictFunction(func(os ...Object) Object {
	if len(os) == 0 {
		return NumArgsError("sub", ">= 1")
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
})

var Mul = NewStrictFunction(func(os ...Object) Object {
	prod := numberType(1)

	for _, o := range os {
		n, ok := o.(numberType)

		if !ok {
			return notNumberError(o)
		}

		prod *= n
	}

	return prod
})

var Div = NewStrictFunction(func(os ...Object) Object {
	if len(os) == 0 {
		return NumArgsError("div", ">= 1")
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
})

var Mod = NewStrictFunction(func(os ...Object) Object {
	if len(os) != 2 {
		return NumArgsError("mod", "2")
	}

	o := os[0]
	n1, ok := o.(numberType)

	if !ok {
		return notNumberError(o)
	}

	o = os[1]
	n2, ok := o.(numberType)

	if !ok {
		return notNumberError(o)
	}

	return NewNumber(math.Mod(float64(n1), float64(n2)))
})

var Pow = NewStrictFunction(func(os ...Object) Object {
	if len(os) != 2 {
		return NumArgsError("pow", "2")
	}

	o := os[0]
	n1, ok := o.(numberType)

	if !ok {
		return notNumberError(o)
	}

	o = os[1]
	n2, ok := o.(numberType)

	if !ok {
		return notNumberError(o)
	}

	return NewNumber(math.Pow(float64(n1), float64(n2)))
})

func notNumberError(o Object) *Thunk {
	return TypeError(o, "Number")
}

// ordered

func (n numberType) less(o ordered) bool {
	return n < o.(numberType)
}
