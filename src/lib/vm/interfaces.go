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
	switch len(ts) {
	case 0:
		return NumArgsError("equal", ">= 1")
	case 1:
		return True
	}

	bs := make([]*Thunk, len(ts)-1)

	for i := range bs {
		bs[i] = equalTwo(ts[0], ts[i+1])
	}

	return And(bs...)
}

func equalTwo(t1, t2 *Thunk) *Thunk {
	ts := []*Thunk{t1, t2}

	for _, t := range ts {
		go t.Eval()
	}

	var es [2]Equalable

	for i, t := range ts {
		o := t.Eval()
		e, ok := o.(Equalable)

		if !ok {
			return notEqualableError(o)
		}

		es[i] = e
	}

	return NewBool(
		reflect.TypeOf(es[0]) == reflect.TypeOf(es[1]) && es[0].Equal(es[1]))
}

func notEqualableError(o Object) *Thunk {
	return TypeError(o, "Equalable")
}

type Addable interface {
	Add(Addable) Addable
}

func Add(ts ...*Thunk) *Thunk {
	if len(ts) == 0 {
		return NumArgsError("add", ">= 1")
	}

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
