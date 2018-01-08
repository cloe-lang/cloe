package core

import (
	"strings"
)

// ListType represents a list of values in the language.
// They can have infinite number of elements inside.
type ListType struct {
	first *Thunk
	rest  *Thunk
}

var (
	emptyList = ListType{nil, nil}

	// EmptyList is a thunk of an empty list.
	EmptyList = Normal(emptyList)
)

// NewList creates a list from its elements.
func NewList(ts ...*Thunk) *Thunk {
	return StrictPrepend(ts, EmptyList)
}

// Prepend prepends multiple elements to a list of the last argument.
var Prepend = NewLazyFunction(
	NewSignature(nil, nil, "elemsAndList", nil, nil, ""),
	prepend)

func prepend(ts ...*Thunk) Value {
	l, err := ts[0].EvalList()

	if err != nil {
		return err
	} else if l.Empty() {
		return NumArgsError("prepend", "> 0")
	}

	if v := ReturnIfEmptyList(l.Rest(), l.First()); v != nil {
		return v
	}

	return cons(l.First(), Normal(prepend(l.Rest())))
}

// StrictPrepend is a strict version of the Prepend function.
func StrictPrepend(ts []*Thunk, l *Thunk) *Thunk {
	for i := len(ts) - 1; i >= 0; i-- {
		l = Normal(cons(ts[i], l))
	}

	return l
}

func cons(t1, t2 *Thunk) ListType {
	return ListType{t1, t2}
}

// First takes the first element in a list.
var First = NewLazyFunction(
	NewSignature([]string{"list"}, nil, "", nil, nil, ""),
	func(ts ...*Thunk) Value {
		l, err := ts[0].EvalList()

		if err != nil {
			return err
		} else if l.Empty() {
			return emptyListError()
		}

		return l.first
	})

// Rest returns a list which has the second to last elements of a given list.
var Rest = NewLazyFunction(
	NewSignature([]string{"list"}, nil, "", nil, nil, ""),
	func(ts ...*Thunk) Value {
		l, err := ts[0].EvalList()

		if err != nil {
			return err
		} else if l.Empty() {
			return emptyListError()
		}

		return PApp(
			NewLazyFunction(
				NewSignature(nil, nil, "", nil, nil, ""),
				func(...*Thunk) Value {
					l, err := l.rest.EvalList()

					if err != nil {
						return err
					}

					return l
				}))
	})

func (l ListType) call(args Arguments) Value {
	return Index.Eval().(callable).call(NewPositionalArguments(Normal(l)).Merge(args))
}

func (l ListType) index(v Value) Value {
	n, err := checkIndex(v)

	if err != nil {
		return err
	}

	for {
		if l.Empty() {
			return OutOfRangeError()
		} else if n == 1 {
			return l.first
		}

		var err Value
		if l, err = l.rest.EvalList(); err != nil {
			return err
		}

		n--
	}
}

func (l ListType) insert(v Value, t *Thunk) Value {
	n, err := checkIndex(v)

	if err != nil {
		return err
	} else if n == 1 {
		return cons(t, Normal(l))
	}

	return cons(l.first, PApp(Insert, l.rest, Normal(n-1), t))
}

func (l ListType) merge(ts ...*Thunk) Value {
	if l.Empty() {
		return PApp(Merge, ts...)
	}

	return cons(l.first, PApp(Merge, append([]*Thunk{l.rest}, ts...)...))
}

func (l ListType) delete(v Value) Value {
	n, err := checkIndex(v)

	if err != nil {
		return err
	}

	es := []*Thunk{}

	for {
		if l.Empty() {
			return OutOfRangeError()
		} else if n == 1 {
			return PApp(Prepend, append(es, l.rest)...)
		}

		es = append(es, l.first)

		var err Value
		if l, err = l.rest.EvalList(); err != nil {
			return err
		}

		n--
	}
}

func checkIndex(v Value) (NumberType, Value) {
	n, ok := v.(NumberType)

	if !ok {
		return 0, NotNumberError(v)
	}

	if !IsInt(n) {
		return 0, NotIntError(n)
	} else if n < 1 {
		return 0, OutOfRangeError()
	}

	return n, nil
}

func (l ListType) toList() Value {
	return l
}

func (l ListType) compare(x comparable) int {
	ll := x.(ListType)

	if l.Empty() && ll.Empty() {
		return 0
	} else if l.Empty() {
		return -1
	} else if ll.Empty() {
		return 1
	}

	c := compare(l.first.Eval(), ll.first.Eval())

	if c == 0 {
		return compare(l.rest.Eval(), ll.rest.Eval())
	}

	return c
}

func (ListType) ordered() {}

func (l ListType) string() Value {
	ss := []string{}

	for !l.Empty() {
		s, err := StrictDump(l.First().Eval())

		if err != nil {
			return err
		}

		ss = append(ss, string(s))

		l, err = l.Rest().EvalList()

		if err != nil {
			return err
		}
	}

	return StringType("[" + strings.Join(ss, " ") + "]")
}

func (l ListType) size() Value {
	n := NumberType(0)

	for !l.Empty() {
		n++

		var err Value
		l, err = l.Rest().EvalList()

		if err != nil {
			return err
		}
	}

	return n
}

func (l ListType) include(elem Value) Value {
	if l.Empty() {
		return False
	}

	b, err := PApp(Equal, l.first, Normal(elem)).EvalBool()

	if err != nil {
		return err
	} else if b {
		return True
	}

	return PApp(Include, l.rest, Normal(elem))
}

// First returns a first element in a list.
func (l ListType) First() *Thunk {
	return l.first
}

// Rest returns elements in a list except the first one.
func (l ListType) Rest() *Thunk {
	return l.rest
}

// Empty returns true if the list is empty.
func (l ListType) Empty() bool {
	return l == ListType{}
}

// ReturnIfEmptyList returns true if a given list is empty, or false otherwise.
func ReturnIfEmptyList(t *Thunk, v Value) Value {
	if l, err := t.EvalList(); err != nil {
		return err
	} else if l.Empty() {
		return v
	}

	return nil
}
