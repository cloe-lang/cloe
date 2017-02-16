package vm

import "fmt"

type ListType struct {
	first *Thunk
	rest  *Thunk
}

var emptyList = ListType{nil, nil}

func NewList(ts ...*Thunk) *Thunk {
	l := Normal(emptyList)

	for i := len(ts) - 1; i >= 0; i-- {
		l = Normal(cons(ts[i], l))
	}

	return l
}

func (l1 ListType) equal(e equalable) Object {
	l2 := e.(ListType)

	if l1 == emptyList || l2 == emptyList {
		return rawBool(l1 == l2)
	}

	for _, t := range []*Thunk{
		// Don't evaluate these parallely for short circuit behavior.
		App(Equal, l1.first, l2.first),
		App(Equal, l1.rest, l2.rest),
	} {
		o := t.Eval()
		b, ok := o.(BoolType)

		if !ok {
			return notBoolError(o)
		} else if !b {
			return False
		}
	}

	return True
}

var Prepend = NewLazyFunction(func(ts ...*Thunk) Object {
	t := ts[0]

	o := ts[1].Eval()
	l, ok := o.(ListType)

	if !ok {
		panic(fmt.Sprintf("Rest arguments must be a list. %v", o))
	}

	ts, err := l.ToThunks()

	if err != nil {
		return err
	}

	for i := len(ts) - 1; i >= 0; i-- {
		t = Normal(cons(ts[i], t))
	}

	return t
})

func cons(t1, t2 *Thunk) ListType {
	return ListType{t1, t2}
}

var First = NewStrictFunction(func(os ...Object) Object {
	if len(os) != 1 {
		return NumArgsError("first", "1")
	}

	o := os[0]
	l, ok := o.(ListType)

	if !ok {
		return notListError(o)
	} else if l == emptyList {
		return emptyListError()
	}

	return l.first
})

var Rest = NewStrictFunction(func(os ...Object) Object {
	if len(os) != 1 {
		return NumArgsError("rest", "1")
	}

	o := os[0]
	l, ok := o.(ListType)

	if !ok {
		return notListError(o)
	} else if l == emptyList {
		return emptyListError()
	}

	return l.rest
})

var Append = NewLazyFunction(appendFunc)

func appendFunc(ts ...*Thunk) Object {
	o := ts[0].Eval()
	l, ok := o.(ListType)

	if !ok {
		return notListError(o)
	}

	if l == emptyList {
		return NewList(ts[1])
	}

	return cons(l.first, App(NewLazyFunction(appendFunc), l.rest, ts[1]))
}

func notListError(o Object) *Thunk {
	return TypeError(o, "List")
}

func emptyListError() *Thunk {
	return ValueError("The list is empty. You cannot apply rest.")
}

func (l ListType) merge(ts ...*Thunk) Object {
	if l == emptyList {
		return App(Merge, ts[0], NewList(ts[1:]...))
	}

	return cons(l.first, App(Merge, l.rest, NewList(ts...)))
}

// ordered

func (l1 ListType) less(ord ordered) bool {
	l2 := ord.(ListType)

	if l2 == emptyList {
		return false
	} else if l1 == emptyList {
		return true
	}

	// Compare firsts

	o1 := l1.first.Eval()
	o2 := l2.first.Eval()

	if less(o1, o2) {
		return true
	} else if less(o2, o1) {
		return false
	}

	// Compare rests

	return less(l1.rest.Eval(), l2.rest.Eval())
}

func (l ListType) ToThunks() ([]*Thunk, *Thunk) {
	ts := make([]*Thunk, 0)

	for l != emptyList {
		ts = append(ts, l.first)

		o := l.rest.Eval()
		var ok bool
		l, ok = o.(ListType)

		if !ok {
			return nil, notListError(o)
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
