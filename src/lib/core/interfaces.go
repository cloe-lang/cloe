package core

import (
	"reflect"
	"strings"
)

// Value represents a value in the language.
// Hackingly, it can be *Thunk so that tail calls are eliminated.
// See also Thunk.Eval().
type Value interface{}

type callable interface {
	call(Arguments) Value // index as function calls
}

// stringable is an interface for something convertable into StringType.
// This should be implemented for all types including error type.
type stringable interface {
	string() Value
}

// ToString converts some value into one of StringType.
var ToString = NewLazyFunction(
	NewSignature(
		[]string{"x"}, nil, "",
		nil, nil, "",
	),
	func(ts ...*Thunk) Value {
		v := ts[0].Eval()
		s, ok := v.(stringable)

		if !ok {
			return TypeError(v, "stringable")
		}

		return s.string()
	})

// Equal returns true when arguments are equal and false otherwise.
// Comparing error values is invalid and it should return an error value.
var Equal = NewStrictFunction(
	NewSignature(
		[]string{"x", "y"}, nil, "",
		nil, nil, "",
	),
	func(ts ...*Thunk) (v Value) {
		defer func() {
			if r := recover(); r != nil {
				v = r
			}
		}()

		return BoolType(compare(ts[0].Eval(), ts[1].Eval()) == 0)
	})

// ordered must be implemented for every type other than error type.
// This interface should not be used in exported functions and exists only to
// make keys for collections in rbt package.
type ordered interface {
	compare(ordered) int // can panic
}

func compare(x1, x2 interface{}) int {
	o1, ok := x1.(ordered)

	if !ok {
		panic(notOrderedError(x1))
	}

	o2, ok := x2.(ordered)

	if !ok {
		panic(notOrderedError(x2))
	}

	if !areSameType(o1, o2) {
		return strings.Compare(reflect.TypeOf(o1).Name(), reflect.TypeOf(o2).Name())
	}

	return o1.compare(o2)
}

func areSameType(x1, x2 interface{}) bool {
	return reflect.TypeOf(x1) == reflect.TypeOf(x2)
}
