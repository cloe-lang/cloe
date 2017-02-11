package vm

type BoolType bool

var True, False = NewBool(true), NewBool(false)

func NewBool(b bool) *Thunk {
	return Normal(rawBool(b))
}

func rawBool(b bool) BoolType {
	return BoolType(b)
}

func (b BoolType) equal(e equalable) Object {
	return rawBool(b == e.(BoolType))
}

var If = NewLazyFunction(func(ts ...*Thunk) Object {
	if len(ts) != 3 {
		return NumArgsError("if", "3")
	}

	o := ts[0].Eval()
	b, ok := o.(BoolType)

	if !ok {
		return notBoolError(o)
	}

	if b {
		return ts[1]
	}

	return ts[2]
})

func notBoolError(o Object) *Thunk {
	return TypeError(o, "Bool")
}

// ordered

func (b BoolType) less(o ordered) bool {
	return bool(!b && o.(BoolType))
}
