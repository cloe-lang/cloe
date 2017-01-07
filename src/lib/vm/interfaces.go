package vm

import "reflect"

type Object interface{}

type Callable interface {
	Call(...*Thunk) *Thunk
}

type Equalable interface {
	Equal(Equalable) bool
}

func Equal(ts ...*Thunk) *Thunk {
	for _, t := range ts {
		go t.Eval()
	}

	o := ts[0].Eval()
	e0, ok := o.(Equalable)

	if !ok {
		return notEqualableError(o)
	}

	for _, t := range ts[1:] {
		if e, ok := t.Eval().(Equalable); !ok {
			return notEqualableError(e)
		} else if reflect.TypeOf(e0) != reflect.TypeOf(e) || !e0.Equal(e) {
			return False
		}
	}

	return True
}

func notEqualableError(o Object) *Thunk {
	return TypeError(o, "Equalable")
}

type Addable interface {
	Add(Addable) Addable
}

func Add(ts ...*Thunk) *Thunk {
	for _, t := range ts {
		go t.Eval()
	}

	o := ts[0].Eval()
	a0, ok := o.(Addable)

	if !ok {
		return TypeError(o, "Addable")
	}

	for _, t := range ts[1:] {
		o := t.Eval()

		if typ := reflect.TypeOf(a0); typ != reflect.TypeOf(o) {
			return TypeError(o, typ.Name())
		}

		a0 = a0.Add(o.(Addable))
	}

	return Normal(a0)
}
