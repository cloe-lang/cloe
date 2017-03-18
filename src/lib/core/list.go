package core

import "strings"

// ListType represents a sequence of values.
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

func (l ListType) equal(e equalable) Value {
	ll := e.(ListType)

	if l == emptyList || ll == emptyList {
		return rawBool(l == ll)
	}

	for _, t := range []*Thunk{
		// Don't evaluate these parallelly for short circuit behavior.
		PApp(Equal, l.first, ll.first),
		PApp(Equal, l.rest, ll.rest),
	} {
		v := t.Eval()
		b, ok := v.(BoolType)

		if !ok {
			return NotBoolError(v)
		} else if !b {
			return False
		}
	}

	return True
}

// Prepend prepends multiple elements to a list of the last argument.
var Prepend = NewLazyFunction(
	NewSignature(
		[]string{}, []OptionalArgument{}, "elemsAndList",
		[]string{}, []OptionalArgument{}, "",
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
var First = NewStrictFunction(
	NewSignature(
		[]string{"list"}, []OptionalArgument{}, "",
		[]string{}, []OptionalArgument{}, "",
	),
	func(vs ...Value) Value {
		v := vs[0]
		l, ok := v.(ListType)

		if !ok {
			return NotListError(v)
		} else if l == emptyList {
			return emptyListError()
		}

		return l.first
	})

// Rest returns a list which has the second to last elements of a given list.
var Rest = NewStrictFunction(
	NewSignature(
		[]string{"list"}, []OptionalArgument{}, "",
		[]string{}, []OptionalArgument{}, "",
	),
	func(vs ...Value) Value {
		v := vs[0]
		l, ok := v.(ListType)

		if !ok {
			return NotListError(v)
		} else if l == emptyList {
			return emptyListError()
		}

		return l.rest
	})

var appendFuncSignature = NewSignature(
	[]string{"list", "elem"}, []OptionalArgument{}, "",
	[]string{}, []OptionalArgument{}, "",
)

// Append appends an element at the end of a given list.
var Append = NewLazyFunction(appendFuncSignature, appendFunc)

func appendFunc(ts ...*Thunk) Value {
	v := ts[0].Eval()
	l, ok := v.(ListType)

	if !ok {
		return NotListError(v)
	}

	if l == emptyList {
		return NewList(ts[1])
	}

	return cons(
		l.first,
		PApp(NewLazyFunction(appendFuncSignature, appendFunc), l.rest, ts[1]),
	)
}

func emptyListError() *Thunk {
	return ValueError("The list is empty. You cannot apply rest.")
}

func (l ListType) call(args Arguments) Value {
	return Index.Eval().(callable).call(NewPositionalArguments(Normal(l)).Merge(args))
}

func (l ListType) index(v Value) (result Value) {
	n, ok := v.(NumberType)

	if !ok {
		return NotNumberError(v)
	}

	v = PApp(isInt, Normal(n)).Eval()
	b, ok := v.(BoolType)

	if !ok {
		return NotBoolError(v)
	} else if !b {
		return NotIntError(n)
	}

	for {
		if l == emptyList {
			return OutOfRangeError()
		} else if n == 0 {
			return l.first
		}

		v = l.rest.Eval()
		l, ok = v.(ListType)

		if !ok {
			return NotListError(v)
		}

		n--
	}
}

func (l ListType) merge(ts ...*Thunk) Value {
	if l == emptyList {
		return PApp(Merge, ts...)
	}

	return cons(l.first, PApp(Merge, append([]*Thunk{l.rest}, ts...)...))
}

func (l ListType) toList() Value {
	return l
}

func (l ListType) less(ord ordered) bool {
	ll := ord.(ListType)

	if ll == emptyList {
		return false
	} else if l == emptyList {
		return true
	}

	// Compare firsts

	o1 := l.first.Eval()
	o2 := ll.first.Eval()

	if less(o1, o2) {
		return true
	} else if less(o2, o1) {
		return false
	}

	// Compare rests

	return less(l.rest.Eval(), ll.rest.Eval())
}

func (l ListType) string() Value {
	vs, err := l.ToValues()

	if err != nil {
		return err.Eval()
	}

	ss := make([]string, len(vs))

	for i, v := range vs {
		if err, ok := v.(ErrorType); ok {
			return err
		}

		v = PApp(ToString, Normal(v)).Eval()
		s, ok := v.(StringType)

		if !ok {
			return NotStringError(v)
		}

		ss[i] = string(s)
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
func (l ListType) ToValues() ([]Value, *Thunk) {
	ts, err := l.ToThunks()

	if err != nil {
		return nil, err
	}

	vs := make([]Value, len(ts))

	for _, t := range ts {
		go t.Eval()
	}

	for i, t := range ts {
		vs[i] = t.Eval()
	}

	return vs, nil
}
