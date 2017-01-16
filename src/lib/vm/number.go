package vm

import (
	"github.com/mediocregopher/seq"
	"math"
)

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

var Sub = NewStrictFunction(func(os ...Object) Object {
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
})

var Mult = NewStrictFunction(func(os ...Object) Object {
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
})

var Div = NewStrictFunction(func(os ...Object) Object {
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
})

var Mod = NewStrictFunction(func(os ...Object) Object {
	if len(os) != 2 {
		return numArgsError("mod", "2")
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

func notNumberError(o Object) *Thunk {
	return typeError(o, "Number")
}

// seq.Setable

func (n numberType) Hash(i uint32) uint32 {
	// TODO: Cast float64 as int64 first by interpreting 64bits and fold it in
	// half.
	return (i + uint32(n)) % seq.ARITY
}

func (n1 numberType) Equal(o interface{}) bool {
	n2, ok := o.(numberType)

	if !ok {
		return false
	}

	return n1 == n2
}
