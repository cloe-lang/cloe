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
	for i, t := range ts {
		o := t.Eval()
		b, ok := o.(Bool)

		if !ok {
			return NewError("%dth argument of and is not Bool but %T.", i, o)
		} else if !bool(b) {
			return False
		}
	}

	return True
}

func Or(ts ...*Thunk) *Thunk { // With short circuit
	for i, t := range ts {
		o := t.Eval()
		b, ok := o.(Bool)

		if !ok {
			return NewError("%dth argument of or is not Bool but %T.", i, o)
		} else if bool(b) {
			return True
		}
	}

	return False
}
