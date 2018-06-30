package core

import "strings"

// ListType represents a list of values in the language.
// They can have infinite number of elements inside.
type ListType struct {
	first Value
	rest  Value
}

// Eval evaluates a value into a WHNF.
func (l *ListType) eval() Value {
	return l
}

var (
	emptyList = ListType{}

	// EmptyList is a thunk of an empty list.
	EmptyList = &emptyList
)

// NewList creates a list from its elements.
func NewList(vs ...Value) Value {
	return StrictPrepend(vs, EmptyList)
}

// Prepend prepends multiple elements to a list of the last argument.
var Prepend = NewLazyFunction(
	NewSignature(nil, "elemsAndList", nil, ""),
	prepend)

func prepend(vs ...Value) Value {
	l, err := EvalList(vs[0])

	if err != nil {
		return err
	}

	if v := ReturnIfEmptyList(l.Rest(), l.First()); v != nil {
		return v
	}

	return cons(l.First(), prepend(l.Rest()))
}

// StrictPrepend is a strict version of the Prepend function.
func StrictPrepend(vs []Value, l Value) Value {
	for i := len(vs) - 1; i >= 0; i-- {
		l = cons(vs[i], l)
	}

	return l
}

func cons(t1, t2 Value) *ListType {
	return &ListType{t1, t2}
}

// First takes the first element in a list.
var First FunctionType

func initFirst() FunctionType {
	return NewLazyFunction(
		NewSignature([]string{"list"}, "", nil, ""),
		func(vs ...Value) Value {
			l, err := EvalList(vs[0])

			if err != nil {
				return err
			}

			return l.First()
		})
}

// Rest returns a list which has the second to last elements of a given list.
var Rest FunctionType

func initRest() FunctionType {
	return NewLazyFunction(
		NewSignature([]string{"list"}, "", nil, ""),
		func(vs ...Value) Value {
			l, err := EvalList(vs[0])

			if err != nil {
				return err
			}

			// TODO: Review this code. Maybe predefine the list check function.
			return PApp(
				NewLazyFunction(
					NewSignature(nil, "", nil, ""),
					func(...Value) Value {
						l, err := EvalList(l.Rest())

						if err != nil {
							return err
						}

						return l
					}))
		})
}

func (l *ListType) assign(i, v Value) Value {
	n, err := checkIndex(i)

	if err != nil {
		return err
	} else if n == 1 {
		return cons(v, l.Rest())
	}

	return cons(l.First(), PApp(Assign, l.Rest(), NewNumber(float64(n-1)), v))
}

func (l *ListType) index(v Value) Value {
	n, err := checkIndex(v)

	if err != nil {
		return err
	}

	for n != 1 {
		var err Value
		if l, err = EvalList(l.Rest()); err != nil {
			return err
		}

		n--
	}

	return l.First()
}

func (l *ListType) insert(i Value, v Value) Value {
	n, err := checkIndex(i)

	if err != nil {
		return err
	} else if n == 1 {
		return cons(v, l)
	}

	return cons(l.First(), PApp(Insert, l.Rest(), NewNumber(float64(n-1)), v))
}

func (l *ListType) merge(vs ...Value) Value {
	if l.Empty() {
		return PApp(Merge, vs...)
	}

	return cons(l.First(), PApp(Merge, append([]Value{l.Rest()}, vs...)...))
}

func (l *ListType) delete(v Value) Value {
	n, err := checkIndex(v)

	if err != nil {
		return err
	}

	es := []Value{}

	for n != 1 {
		es = append(es, l.First())

		var err Value
		if l, err = EvalList(l.Rest()); err != nil {
			return err
		}

		n--
	}

	return StrictPrepend(es, l.Rest())
}

func checkIndex(v Value) (NumberType, Value) {
	n, err := EvalNumber(v)

	if err != nil {
		return 0, err
	}

	if !IsInt(n) {
		return 0, NotIntError(n)
	} else if n < 1 {
		return 0, OutOfRangeError()
	}

	return n, nil
}

func (l *ListType) toList() Value {
	return l
}

func (l *ListType) compare(x comparable) int {
	ll := x.(*ListType)

	if l.Empty() && ll.Empty() {
		return 0
	} else if l.Empty() {
		return -1
	} else if ll.Empty() {
		return 1
	}

	c := compare(EvalPure(l.First()), EvalPure(ll.First()))

	if c == 0 {
		return compare(EvalPure(l.Rest()), EvalPure(ll.Rest()))
	}

	return c
}

func (*ListType) ordered() {}

func (l *ListType) string() Value {
	ss := []string{}

	for !l.Empty() {
		s, err := StrictDump(EvalPure(l.First()))

		if err != nil {
			return err
		}

		ss = append(ss, string(s))

		if l, err = EvalList(l.Rest()); err != nil {
			return err
		}
	}

	return NewString("[" + strings.Join(ss, " ") + "]")
}

func (l *ListType) size() Value {
	n := NewNumber(0)

	for !l.Empty() {
		*n++

		var err Value
		if l, err = EvalList(l.Rest()); err != nil {
			return err
		}
	}

	return n
}

func (l *ListType) include(elem Value) Value {
	if l.Empty() {
		return False
	}

	b, err := EvalBoolean(PApp(Equal, l.First(), elem))

	if err != nil {
		return err
	} else if b {
		return True
	}

	return PApp(Include, l.Rest(), elem)
}

// First returns a first element in a list.
func (l *ListType) First() Value {
	if l.Empty() {
		return emptyListError()
	}

	return l.first
}

// Rest returns elements in a list except the first one.
func (l *ListType) Rest() Value {
	if l.Empty() {
		return emptyListError()
	}

	return l.rest
}

// Empty returns true if the list is empty.
func (l *ListType) Empty() bool {
	return *l == emptyList
}

// ReturnIfEmptyList returns true if a given list is empty, or false otherwise.
func ReturnIfEmptyList(t Value, v Value) Value {
	if l, err := EvalList(t); err != nil {
		return err
	} else if l.Empty() {
		return v
	}

	return nil
}
