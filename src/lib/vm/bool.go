package vm

type boolType bool

var True, False = NewBool(true), NewBool(false)

func NewBool(b bool) *Thunk {
	return Normal(rawBool(b))
}

func rawBool(b bool) boolType {
	return boolType(b)
}

func (b boolType) Equal(e Equalable) Object {
	return rawBool(b == e.(boolType))
}

var If = NewLazyFunction(ifFunc)

func ifFunc(ts ...*Thunk) Object {
	if len(ts) != 3 {
		return NumArgsError("if", "3")
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
}

func notBoolError(o Object) *Thunk {
	return TypeError(o, "Bool")
}
