package vm

import "reflect"

type Object interface{}

type Callable interface {
	Call(...*Thunk) Object
}

type Equalable interface {
	Equal(Equalable) Object
}

var Equal = NewStrictFunction(equal)

func equal(os ...Object) Object {
	if len(os) != 2 {
		return numArgsError("equal", "2")
	}

	var es [2]Equalable

	for i, o := range os {
		e, ok := o.(Equalable)

		if !ok {
			return typeError(o, "Equalable")
		}

		es[i] = e
	}

	if reflect.TypeOf(es[0]) != reflect.TypeOf(es[1]) {
		return False
	}

	return es[0].Equal(es[1])
}

type Addable interface {
	Add(Addable) Addable
}

// Add should be lazy for list concatenation and dictionary merging.
// THE SENTENCE ABOVE IS WRONG because types of objects must be known to sum up
// them.
var Add = NewStrictFunction(add)

func add(os ...Object) Object {
	if len(os) == 0 {
		return numArgsError("add", ">= 1")
	}

	o := os[0]
	a0, ok := o.(Addable)

	if !ok {
		return typeError(o, "Addable")
	}

	for _, o := range os[1:] {
		if typ := reflect.TypeOf(a0); typ != reflect.TypeOf(o) {
			return typeError(o, typ.Name())
		}

		a0 = a0.Add(o.(Addable))
	}

	return a0
}
