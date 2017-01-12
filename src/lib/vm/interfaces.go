package vm

import "reflect"

type Object interface{}

type callable interface {
	call(...*Thunk) Object
}

type equalable interface {
	equal(equalable) Object
}

var Equal = NewStrictFunction(equal)

func equal(os ...Object) Object {
	if len(os) != 2 {
		return numArgsError("equal", "2")
	}

	var es [2]equalable

	for i, o := range os {
		e, ok := o.(equalable)

		if !ok {
			return typeError(o, "Equalable")
		}

		es[i] = e
	}

	if reflect.TypeOf(es[0]) != reflect.TypeOf(es[1]) {
		return False
	}

	return es[0].equal(es[1])
}

type addable interface {
	add(addable) addable
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
	a0, ok := o.(addable)

	if !ok {
		return typeError(o, "addable")
	}

	for _, o := range os[1:] {
		if typ := reflect.TypeOf(a0); typ != reflect.TypeOf(o) {
			return typeError(o, typ.Name())
		}

		a0 = a0.add(o.(addable))
	}

	return a0
}
