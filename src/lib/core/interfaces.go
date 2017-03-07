package core

import (
	"github.com/raviqqe/tisp/src/lib/util"
	"reflect"
)

// Object represents an object in the language.
// Hackingly, it can be *Thunk so that tail calls are eliminated.
// See also Thunk.Eval().
type Object interface{}

type callable interface {
	call(Arguments) Object
}

// equalable must be implemented for every type other than error type.
type equalable interface {
	equal(equalable) Object
}

// Equal returns true when arguments are equal and false otherwise.
// Comparing error objects is invalid and it should return an error object.
var Equal = NewStrictFunction(
	NewSignature(
		[]string{"x", "y"}, []OptionalArgument{}, "",
		[]string{}, []OptionalArgument{}, "",
	),
	func(os ...Object) Object {
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

var ToList = NewStrictFunction(
	NewSignature(
		[]string{"listLike"}, []OptionalArgument{}, "",
		[]string{}, []OptionalArgument{}, "",
	),
	func(os ...Object) Object {
		o := os[0]
		l, ok := o.(listable)

		if !ok {
			return TypeError(o, "Listable")
		}

		return l.toList()
	})

// ordered must be implemented for every type other than error type.
// This interface should not be used in exported functions and exists only to
// make keys for collections in rbt package.
type ordered interface {
	less(ordered) bool // can panic
}

func less(x1, x2 interface{}) bool {
	o1, ok := x1.(ordered)

	if !ok {
		panic(notOrderedError(x1))
	}

	o2, ok := x2.(ordered)

	if !ok {
		panic(notOrderedError(x2))
	}

	if !areSameType(o1, o2) {
		return reflect.TypeOf(o1).Name() < reflect.TypeOf(o2).Name()
	}

	return o1.less(o2)
}

func areSameType(x1, x2 interface{}) bool {
	return reflect.TypeOf(x1) == reflect.TypeOf(x2)
}

type mergable interface {
	merge(ts ...*Thunk) Object
}

var Merge = NewLazyFunction(
	NewSignature(
		[]string{"x"}, []OptionalArgument{}, "ys",
		[]string{}, []OptionalArgument{}, "",
	),
	func(ts ...*Thunk) Object {
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

type deletable interface {
	delete(Object) (deletable, Object)
}

var Delete = NewStrictFunction(
	NewSignature(
		[]string{"collection"}, []OptionalArgument{}, "elems",
		[]string{}, []OptionalArgument{}, "",
	),
	func(os ...Object) Object {
		d, ok := os[0].(deletable)

		if !ok {
			return TypeError(os[0], "Deletable")
		}

		l, ok := os[1].(ListType)

		if !ok {
			return notListError(os[1])
		}

		os, err := l.ToObjects()

		if err != nil {
			return err
		}

		for _, o := range os {
			if _, ok := o.(ErrorType); ok {
				return o
			}

			var err Object
			d, err = d.delete(o)

			if err != nil {
				return err
			}
		}

		return d
	})

// stringable is an interface for something convertable into StringType.
// This should be implemented for all types including error type.
type stringable interface {
	string() Object
}

// ToString converts some object into one of StringType.
var ToString = NewStrictFunction(
	NewSignature(
		[]string{"x"}, []OptionalArgument{}, "",
		[]string{}, []OptionalArgument{}, "",
	),
	func(os ...Object) Object {
		s, ok := os[0].(stringable)

		if !ok {
			util.Fail("%#v is not stringable.", os[0])
		}

		return s.string()
	})

// TODO: Create collection interface integrating some existing interfaces with
// methods of index, insert, merge, delete, size (or len?), include and toList.
// It should be implemented by StringType, ListType, DictionaryType and SetType.
