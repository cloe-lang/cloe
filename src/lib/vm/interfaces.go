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

	if !areSameType(es[0], es[1]) {
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

// This interface should not be used in exported Functions and exists only to
// make keys of DictionaryType and elements of setType ordered.
type ordered interface {
	less(ordered) bool // can panic
}

func less(x1, x2 interface{}) bool {
	if !areSameType(x1, x2) {
		return reflect.TypeOf(x1).Name() < reflect.TypeOf(x2).Name()
	}

	o1, ok := x1.(ordered)

	if !ok {
		panic(notOrderedError(x1))
	}

	o2, ok := x2.(ordered)

	if !ok {
		panic(notOrderedError(x2))
	}

	return o1.less(o2)
}

func areSameType(x1, x2 interface{}) bool {
	return reflect.TypeOf(x1) == reflect.TypeOf(x2)
}

type mergable interface {
	merge(ts ...*Thunk) Object
}

var Merge = NewLazyFunction(func(ts ...*Thunk) Object {
	o := ts[0].Eval()
	m, ok := o.(mergable)

	if !ok {
		return TypeError(o, "Mergable")
	}

	o = ts[1].Eval()
	l, ok := o.(ListType)

	if !ok {
		return notListError(o)
	}

	ts, err := l.ToThunks()

	if err != nil {
		return err
	}

	if len(ts) == 0 {
		return m
	}

	return m.merge(ts...)
})
