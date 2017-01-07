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

	e0, ok := ts[0].Eval().(Equalable)

	if !ok {
		return notEqualableError(e0)
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
