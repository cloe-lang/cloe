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
	NewSignature([]string{"arg"}, nil, "", nil, nil, ""),
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

type ordered interface {
	comparable
	ordered()
}

// Equal checks if all arguments are equal or not.
var Equal = NewLazyFunction(
	NewSignature(nil, nil, "args", nil, nil, ""),
	func(ts ...*Thunk) (v Value) {
		defer func() {
			if r := recover(); r != nil {
				v = r
			}
		}()

		t := ts[0]
		l, err := t.EvalList()

		if err != nil {
			return err
		} else if l == emptyList {
			return True
		}

		e := l.first.Eval()

		for {
			l, err = l.rest.EvalList()

			if err != nil {
				return err
			} else if l == emptyList {
				return True
			}

			if compare(e, l.first.Eval()) != 0 {
				return False
			}
		}
	})

// Compare compares 2 values and returns -1 when x < y, 0 when x = y, and 1 when x > y.
var Compare = NewStrictFunction(
	NewSignature([]string{"left", "right"}, nil, "", nil, nil, ""),
	rawCompare)

func rawCompare(ts ...*Thunk) Value {
	v := ts[0].Eval()
	o1, ok := v.(ordered)

	if !ok {
		return NotOrderedError(v)
	}

	v = ts[1].Eval()
	o2, ok := v.(ordered)

	if !ok {
		return NotOrderedError(v)
	}

	if reflect.TypeOf(o1) != reflect.TypeOf(o2) {
		s, err := PApp(TypeOf, ts[1]).EvalString()

		if err != nil {
			return err
		}

		return TypeError(o1, string(s))
	}

	if l1, ok := o1.(ListType); ok {
		return compareListsAsOrdered(l1, o2.(ListType))
	}

	c := o1.compare(o2)
	if c < 0 {
		return NewNumber(-1)
	} else if c > 0 {
		return NewNumber(1)
	}

	return NewNumber(0)
}
