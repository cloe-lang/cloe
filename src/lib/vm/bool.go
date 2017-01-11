package vm

type Bool bool

var True, False = NewBool(true), NewBool(false)

func NewBool(b bool) Bool {
	return Bool(b)
}

func BoolThunk(b bool) *Thunk {
	return Normal(NewBool(b))
}

func (b Bool) Equal(e Equalable) Object {
	return NewBool(b == e.(Bool))
}

var If = NewLazyFunction(ifFunc)

func ifFunc(ts ...*Thunk) Object {
	if len(ts) != 3 {
		return NumArgsError("if", "3")
	}

	o := ts[0].Eval()
	b, ok := o.(Bool)

	if !ok {
		return notBoolError(o)
	}

	if b {
		return ts[1]
	}

	return ts[2]
}

func notBoolError(o Object) Error {
	return TypeError(o, "Bool")
}
