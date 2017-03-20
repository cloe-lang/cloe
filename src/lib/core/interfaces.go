package core

import "reflect"

// Value represents a value in the language.
// Hackingly, it can be *Thunk so that tail calls are eliminated.
// See also Thunk.Eval().
type Value interface{}

type callable interface {
	call(Arguments) Value
}

// stringable is an interface for something convertable into StringType.
// This should be implemented for all types including error type.
type stringable interface {
	string() Value
}

// ToString converts some value into one of StringType.
var ToString = NewStrictFunction(
	NewSignature(
		[]string{"x"}, []OptionalArgument{}, "",
		[]string{}, []OptionalArgument{}, "",
	),
	func(vs ...Value) Value {
		s, ok := vs[0].(stringable)

		if !ok {
			return TypeError(vs[0], "stringable")
		}

		return s.string()
	})

// equalable must be implemented for every type other than error type.
type equalable interface {
	equal(equalable) Value
}

// Equal returns true when arguments are equal and false otherwise.
// Comparing error values is invalid and it should return an error value.
var Equal = NewStrictFunction(
	NewSignature(
		[]string{"x", "y"}, []OptionalArgument{}, "",
		[]string{}, []OptionalArgument{}, "",
	),
	func(vs ...Value) Value {
		var es [2]equalable

		for i, v := range vs {
			e, ok := v.(equalable)

			if !ok {
				return TypeError(v, "equalable")
			}

			es[i] = e
		}

		if !areSameType(es[0], es[1]) {
			return False
		}

		return es[0].equal(es[1])
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
