package vm

type Bool bool

var True = NewBool(true)
var False = NewBool(false)

func NewBool(b bool) *Thunk {
	return Normal(Bool(b))
}

func (b Bool) Equal(e Equalable) bool {
	return b == e.(Bool)
}

func And(ts ...*Thunk) *Thunk { // With short circuit
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

func Or(ts ...*Thunk) *Thunk { // With short circuit
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

func notBoolError(o Object) *Thunk {
	return TypeError(o, "Bool")
}
