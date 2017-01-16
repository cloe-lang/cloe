package vm

import "github.com/mediocregopher/seq"

type boolType bool

var True, False = NewBool(true), NewBool(false)

func NewBool(b bool) *Thunk {
	return Normal(rawBool(b))
}

func rawBool(b bool) boolType {
	return boolType(b)
}

func (b boolType) equal(e equalable) Object {
	return rawBool(b == e.(boolType))
}

var If = NewLazyFunction(func(ts ...*Thunk) Object {
	if len(ts) != 3 {
		return numArgsError("if", "3")
	}

	o := ts[0].Eval()
	b, ok := o.(boolType)

	if !ok {
		return notBoolError(o)
	}

	if b {
		return ts[1]
	}

	return ts[2]
})

func notBoolError(o Object) *Thunk {
	return typeError(o, "Bool")
}

// seq.Setable

func (b boolType) Hash(i uint32) uint32 {
	var j uint32

	if b {
		j = 1
	} else {
		j = 0
	}

	return (i + j) % seq.ARITY
}

func (b1 boolType) Equal(o interface{}) bool {
	b2, ok := o.(boolType)

	if !ok {
		return false
	}

	return b1 == b2
}
