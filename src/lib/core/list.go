package core

import (
	"strings"

	"github.com/coel-lang/coel/src/lib/systemt"
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
	l := Normal(emptyList)

	for i := len(ts) - 1; i >= 0; i-- {
		l = Normal(cons(ts[i], l))
	}

	return l
}

// Prepend prepends multiple elements to a list of the last argument.
var Prepend = NewLazyFunction(
	NewSignature(
		nil, nil, "elemsAndList",
		nil, nil, "",
	),
	func(ts ...*Thunk) Value {
		v := ts[0].Eval()
		ts, err := v.(ListType).ToThunks()

		if err != nil {
			return err
		} else if len(ts) == 0 {
			return NumArgsError("prepend", "> 0")
		}

		last := len(ts) - 1
		t := ts[last]

		for i := last - 1; i >= 0; i-- {
			t = Normal(cons(ts[i], t))
		}

		return t
	})

func cons(t1, t2 *Thunk) ListType {
	return ListType{t1, t2}
}

// First takes the first element in a list.
var First = NewLazyFunction(
	NewSignature(
		[]string{"list"}, nil, "",
		nil, nil, "",
	),
	func(ts ...*Thunk) Value {
		v := ts[0].Eval()
		l, ok := v.(ListType)

		if !ok {
			return NotListError(v)
		} else if l == emptyList {
			return emptyListError()
		}

		return l.first
	})

// Rest returns a list which has the second to last elements of a given list.
var Rest = NewLazyFunction(
	NewSignature(
		[]string{"list"}, nil, "",
		nil, nil, "",
	),
	func(ts ...*Thunk) Value {
		v := ts[0].Eval()
		l, ok := v.(ListType)

		if !ok {
			return NotListError(v)
		} else if l == emptyList {
			return emptyListError()
		}

		return PApp(
			NewLazyFunction(
				NewSignature(nil, nil, "", nil, nil, ""),
				func(...*Thunk) Value {
					v := l.rest.Eval()
					l, ok := v.(ListType)

					if !ok {
						return NotListError(v)
					}

					return l
				}))
	})

var appendFuncSignature = NewSignature([]string{"list"}, nil, "elems", nil, nil, "")

func appendFunc(ts ...*Thunk) Value {
	return PApp(Merge, ts...)
}

// Append appends an element at the end of a given list.
var Append = NewLazyFunction(appendFuncSignature, appendFunc)

func (l ListType) call(args Arguments) Value {
	return Index.Eval().(callable).call(NewPositionalArguments(Normal(l)).Merge(args))
}

func (l ListType) index(v Value) Value {
	n, err := checkIndex(v)

	if err != nil {
		return err
	}

	for {
		if l == emptyList {
			return OutOfRangeError()
		} else if n == 0 {
			return l.first
		}

		v = l.rest.Eval()
		var ok bool
		l, ok = v.(ListType)

		if !ok {
			return NotListError(v)
		}

		n--
	}
}

func (l ListType) insert(ts ...*Thunk) Value {
	if len(ts) != 2 {
		return NumArgsError("insert", "3 if a collection is a list")
	}

	v := ts[0].Eval()
	n, err := checkIndex(v)

	if err != nil {
		return err
	} else if n == 0 {
		return PApp(Prepend, ts[1], Normal(l))
	}

	return PApp(Prepend, l.first, PApp(Insert, l.rest, Normal(n-1), ts[1]))
}

func (l ListType) merge(ts ...*Thunk) Value {
	if l == emptyList {
		return PApp(Merge, ts...)
	}

	return cons(l.first, PApp(Merge, append([]*Thunk{l.rest}, ts...)...))
}

func (l ListType) delete(v Value) Value {
	n, err := checkIndex(v)

	if err != nil {
		return err
	}

	elems := make([]*Thunk, 0)

	for {
		if l == emptyList {
			return OutOfRangeError()
		} else if n == 0 {
			return PApp(Merge, NewList(elems...), l.rest)
		}

		elems = append(elems, l.first)
		v = l.rest.Eval()
		var ok bool
		l, ok = v.(ListType)

		if !ok {
			return NotListError(v)
		}

		n--
	}
}

func checkIndex(v Value) (NumberType, Value) {
	n, ok := v.(NumberType)

	if !ok {
		return 0, NotNumberError(v)
	}

	v = PApp(isInt, Normal(n)).Eval()
	b, ok := v.(BoolType)

	if !ok {
		return 0, NotBoolError(v)
	} else if !b {
		return 0, NotIntError(n)
	} else if n < 0 {
		return 0, OutOfRangeError()
	}

	return n, nil
}

func (l ListType) toList() Value {
	return l
}

func (l ListType) compare(x comparable) int {
	ll := x.(ListType)

	if l == emptyList && ll == emptyList {
		return 0
	} else if l == emptyList {
		return -1
	} else if ll == emptyList {
		return 1
	}

	c := compare(l.first.Eval(), ll.first.Eval())

	if c == 0 {
		return compare(l.rest.Eval(), ll.rest.Eval())
	}

	return c
}

func (ListType) ordered() {}

func compareListsAsOrdered(l, ll ListType) Value {
	if l == emptyList && ll == emptyList {
		return NewNumber(0)
	} else if l == emptyList {
		return NewNumber(-1)
	} else if ll == emptyList {
		return NewNumber(1)
	}

	v := ensureNormal(rawCompare(l.first, ll.first))
	n, ok := v.(NumberType)

	if !ok {
		return NotNumberError(v)
	} else if n == 0 {
		return rawCompare(l.rest, ll.rest)
	}

	return n
}

func (l ListType) string() Value {
	ts, err := l.ToValues()

	if err != nil {
		return err.Eval()
	}

	ss := make([]string, 0, len(ts))

	for _, t := range ts {
		v := PApp(Dump, t).Eval()
		s, ok := v.(StringType)

		if !ok {
			return NotStringError(v)
		}

		ss = append(ss, string(s))
	}

	return StringType("[" + strings.Join(ss, " ") + "]")
}

// ToThunks converts a list into a slice of its elements as thunks.
func (l ListType) ToThunks() ([]*Thunk, *Thunk) {
	ts := make([]*Thunk, 0)

	for l != emptyList {
		ts = append(ts, l.first)

		v := l.rest.Eval()
		var ok bool
		l, ok = v.(ListType)

		if !ok {
			return nil, NotListError(v)
		}
	}

	return ts, nil
}

// ToValues converts a list into a slice of its elements as values.
func (l ListType) ToValues() ([]*Thunk, *Thunk) {
	ts, err := l.ToThunks()

	if err != nil {
		return nil, err
	}

	for _, t := range ts {
		tt := t
		systemt.Daemonize(func() { tt.Eval() })
	}

	return ts, nil
}

func (l ListType) size() Value {
	ts, err := l.ToThunks()

	if err != nil {
		return err
	}

	return NumberType(len(ts))
}

func (l ListType) include(elem Value) Value {
	if l == emptyList {
		return False
	}

	v := PApp(Equal, l.first, Normal(elem)).Eval()
	b, ok := v.(BoolType)

	if !ok {
		return NotBoolError(v)
	} else if b {
		return True
	}

	return PApp(Include, l.rest, Normal(elem))
}
