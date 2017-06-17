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

// comparable must be implemented for every type other than error type.
// This interface should not be used in exported functions and exists only to
// make keys for collections in rbt package.
type comparable interface {
	compare(comparable) int // can panic
}

func compare(x1, x2 interface{}) int {
	o1, ok := x1.(comparable)

	if !ok {
		panic(notComparableError(x1))
	}

	o2, ok := x2.(comparable)

	if !ok {
		panic(notComparableError(x2))
	}

	if reflect.TypeOf(o1) != reflect.TypeOf(o2) {
		return strings.Compare(reflect.TypeOf(o1).Name(), reflect.TypeOf(o2).Name())
	}

	return o1.compare(o2)
}

// Compare compares 2 values and returns -1 when x < y, 0 when x = y, and 1 when x > y.
var Compare = NewStrictFunction(
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

		c := compare(ts[0].Eval(), ts[1].Eval())
		if c < 0 {
			return NewNumber(-1)
		} else if c > 0 {
			return NewNumber(1)
		}

		return NewNumber(0)
	})
