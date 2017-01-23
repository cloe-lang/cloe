package vm

import "reflect"

type Object interface{}

type callable interface {
	call(...*Thunk) Object
}

type equalable interface {
	equal(equalable) Object
}

var Equal = NewStrictFunction(func(os ...Object) Object {
	if len(os) != 2 {
		return NumArgsError("equal", "2")
	}

	var es [2]equalable

	for i, o := range os {
		e, ok := o.(equalable)

		if !ok {
			return TypeError(o, "Equalable")
		}

		es[i] = e
	}

	if reflect.TypeOf(es[0]) != reflect.TypeOf(es[1]) {
		return False
	}

	return es[0].equal(es[1])
})

type listable interface {
	toList() Object
}

var ToList = NewStrictFunction(func(os ...Object) Object {
	if len(os) != 1 {
		return NumArgsError("toList", "1")
	}

	o := os[0]
	l, ok := o.(listable)

	if !ok {
		return TypeError(o, "Listable")
	}

	return l.toList()
})
