package vm

type Bool bool

var True, False = NewBool(true), NewBool(false)

func NewBool(b bool) *Thunk {
	return Normal(Bool(b))
}

func (b Bool) Equal(e Equalable) bool {
	return b == e.(Bool)
}

func And(ts ...*Thunk) *Thunk { // with short circuit
	for _, t := range ts {
		o := t.Eval()
		b, ok := o.(Bool)

		if !ok {
			return notBoolError(o)
		} else if !bool(b) {
			return False
		}
	}

	return True
}

func Or(ts ...*Thunk) *Thunk { // with short circuit
	for _, t := range ts {
		o := t.Eval()
		b, ok := o.(Bool)

		if !ok {
			return notBoolError(o)
		} else if bool(b) {
			return True
		}
	}

	return False
}

func Not(ts ...*Thunk) *Thunk {
	if len(ts) != 1 {
		return NumArgsError("not", "1")
	}

	o := ts[0].Eval()
	b, ok := o.(Bool)

	if !ok {
		return notBoolError(o)
	}

	return NewBool(!bool(b))
}

func notBoolError(o Object) *Thunk {
	return TypeError(o, "Bool")
}
