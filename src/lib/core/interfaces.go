package core

import (
	"reflect"
)

type callable interface {
	Value
	call(Arguments) Value // index as function calls
}

// stringable is an interface for something convertable into StringType.
// This should be implemented for all types including error type.
type stringable interface {
	Value
	string() Value
}

// ToString converts some value into one of StringType.
var ToString = NewLazyFunction(
	NewSignature([]string{"arg"}, "", nil, ""),
	func(vs ...Value) Value {
		v := EvalPure(vs[0])
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
	Value
	compare(comparable) int // can panic
}

func compare(x1, x2 interface{}) int {
	v1, v2 := x1.(Value), x2.(Value)
	o1, ok := v1.(comparable)

	if !ok {
		panic(notComparableError(v1))
	}

	o2, ok := v2.(comparable)

	if !ok {
		panic(notComparableError(v2))
	}

	i1 := comparableID(v1)
	i2 := comparableID(v2)

	if i1 == i2 {
		return o1.compare(o2)
	} else if i1 < i2 {
		return -1
	}

	return 1
}

func comparableID(v Value) byte {
	switch v.(type) {
	case *BooleanType:
		return 0
	case *DictionaryType:
		return 1
	case *ListType:
		return 2
	case NilType:
		return 3
	case *NumberType:
		return 4
	case StringType:
		return 5
	}

	panic("Unreachable")
}

type ordered interface {
	comparable
	ordered()
}

// Equal checks if all arguments are equal or not.
var Equal = NewLazyFunction(
	NewSignature(nil, "args", nil, ""),
	func(vs ...Value) (v Value) {
		defer func() {
			if r := recover(); r != nil {
				v = r.(Value)
			}
		}()

		l, err := EvalList(vs[0])

		if err != nil {
			return err
		} else if l.Empty() {
			return True
		}

		e := EvalPure(l.First())

		for {
			l, err = EvalList(l.Rest())

			if err != nil {
				return err
			} else if l.Empty() {
				return True
			}

			if compare(e, EvalPure(l.First())) != 0 {
				return False
			}
		}
	})

// Compare compares 2 values and returns -1 when x < y, 0 when x = y, and 1 when x > y.
var Compare = NewStrictFunction(
	NewSignature([]string{"left", "right"}, "", nil, ""),
	compareAsOrdered)

func compareAsOrdered(vs ...Value) Value {
	v := EvalPure(vs[0])
	o1, ok := v.(ordered)

	if !ok {
		return NotOrderedError(v)
	}

	v = EvalPure(vs[1])
	o2, ok := v.(ordered)

	if !ok {
		return NotOrderedError(v)
	}

	if reflect.TypeOf(o1) != reflect.TypeOf(o2) {
		s, err := EvalString(PApp(TypeOf, vs[1]))

		if err != nil {
			return err
		}

		return TypeError(o1, string(s))
	}

	if l1, ok := o1.(*ListType); ok {
		return compareListsAsOrdered(l1, o2.(*ListType))
	}

	c := o1.compare(o2)

	if c < 0 {
		return NewNumber(-1)
	} else if c > 0 {
		return NewNumber(1)
	}

	return NewNumber(0)
}

func compareListsAsOrdered(l, ll *ListType) Value {
	if l.Empty() && ll.Empty() {
		return NewNumber(0)
	} else if l.Empty() {
		return NewNumber(-1)
	} else if ll.Empty() {
		return NewNumber(1)
	}

	v := compareAsOrdered(l.First(), ll.First())
	n, err := EvalNumber(v)

	if err != nil {
		return err
	} else if n == 0 {
		return compareAsOrdered(l.Rest(), ll.Rest())
	}

	return &n
}

// IsOrdered checks if a value is ordered or not.
var IsOrdered = NewLazyFunction(
	NewSignature([]string{"arg"}, "", nil, ""),
	isOrdered)

func isOrdered(vs ...Value) Value {
	switch x := EvalPure(vs[0]).(type) {
	case *ErrorType:
		return x
	case *ListType:
		for !x.Empty() {
			b, err := EvalBoolean(isOrdered(x.First()))

			if err != nil {
				return err
			} else if !b {
				return False
			}

			if x, err = EvalList(x.Rest()); err != nil {
				return err
			}
		}

		return True
	default:
		_, ok := x.(ordered)
		return NewBoolean(ok)
	}
}
