package core

import "strings"

type ListType struct {
	first *Thunk
	rest  *Thunk
}

var (
	emptyList = ListType{nil, nil}
	EmptyList = Normal(emptyList)
)

func NewList(ts ...*Thunk) *Thunk {
	l := Normal(emptyList)

	for i := len(ts) - 1; i >= 0; i-- {
		l = Normal(cons(ts[i], l))
	}

	return l
}

func (l ListType) equal(e equalable) Object {
	ll := e.(ListType)

	if l == emptyList || ll == emptyList {
		return rawBool(l == ll)
	}

	for _, t := range []*Thunk{
		// Don't evaluate these parallelly for short circuit behavior.
		PApp(Equal, l.first, ll.first),
		PApp(Equal, l.rest, ll.rest),
	} {
		o := t.Eval()
		b, ok := o.(BoolType)

		if !ok {
			return NotBoolError(o)
		} else if !b {
			return False
		}
	}

	return True
}

var Prepend = NewLazyFunction(
	NewSignature(
		[]string{}, []OptionalArgument{}, "elemsAndList",
		[]string{}, []OptionalArgument{}, "",
	),
	func(ts ...*Thunk) Object {
		o := ts[0].Eval()
		ts, err := o.(ListType).ToThunks()

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

var First = NewStrictFunction(
	NewSignature(
		[]string{"list"}, []OptionalArgument{}, "",
		[]string{}, []OptionalArgument{}, "",
	),
	func(os ...Object) Object {
		o := os[0]
		l, ok := o.(ListType)

		if !ok {
			return NotListError(o)
		} else if l == emptyList {
			return emptyListError()
		}

		return l.first
	})

var Rest = NewStrictFunction(
	NewSignature(
		[]string{"list"}, []OptionalArgument{}, "",
		[]string{}, []OptionalArgument{}, "",
	),
	func(os ...Object) Object {
		o := os[0]
		l, ok := o.(ListType)

		if !ok {
			return NotListError(o)
		} else if l == emptyList {
			return emptyListError()
		}

		return l.rest
	})

var appendFuncSignature = NewSignature(
	[]string{"list", "elem"}, []OptionalArgument{}, "",
	[]string{}, []OptionalArgument{}, "",
)

var Append = NewLazyFunction(appendFuncSignature, appendFunc)

func appendFunc(ts ...*Thunk) Object {
	o := ts[0].Eval()
	l, ok := o.(ListType)

	if !ok {
		return NotListError(o)
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

func (l ListType) merge(ts ...*Thunk) Object {
	if l == emptyList {
		return PApp(Merge, ts...)
	}

	return cons(l.first, PApp(Merge, append([]*Thunk{l.rest}, ts...)...))
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

func (l ListType) string() Object {
	os, err := l.ToObjects()

	if err != nil {
		return err.Eval()
	}

	ss := make([]string, len(os))

	for i, o := range os {
		if err, ok := o.(ErrorType); ok {
			return err
		}

		o = PApp(ToString, Normal(o)).Eval()
		s, ok := o.(StringType)

		if !ok {
			return NotStringError(o)
		}

		ss[i] = string(s)
	}

	return StringType("[" + strings.Join(ss, " ") + "]")
}

func (l ListType) ToThunks() ([]*Thunk, *Thunk) {
	ts := make([]*Thunk, 0)

	for l != emptyList {
		ts = append(ts, l.first)

		o := l.rest.Eval()
		var ok bool
		l, ok = o.(ListType)

		if !ok {
			return nil, NotListError(o)
		}
	}

	return ts, nil
}

func (l ListType) ToObjects() ([]Object, *Thunk) {
	ts, err := l.ToThunks()

	if err != nil {
		return nil, err
	}

	os := make([]Object, len(ts))

	for _, t := range ts {
		go t.Eval()
	}

	for i, t := range ts {
		os[i] = t.Eval()
	}

	return os, nil
}
