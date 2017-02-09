package vm

type listType struct {
	first *Thunk
	rest  *Thunk
}

var emptyList = listType{nil, nil}

func NewList(ts ...*Thunk) *Thunk {
	return App(Prepend, append(ts, Normal(emptyList))...)
}

func (l1 listType) equal(e equalable) Object {
	l2 := e.(listType)

	if l1 == emptyList || l2 == emptyList {
		return rawBool(l1 == l2)
	}

	for _, t := range []*Thunk{
		// Don't evaluate these parallely for short circuit behavior.
		App(Equal, l1.first, l2.first),
		App(Equal, l1.rest, l2.rest),
	} {
		o := t.Eval()
		b, ok := o.(boolType)

		if !ok {
			return notBoolError(o)
		} else if !b {
			return False
		}
	}

	return True
}

var Prepend = NewLazyFunction(func(ts ...*Thunk) Object {
	switch len(ts) {
	case 0:
		return NumArgsError("prepend", "> 1")
	case 1:
		return ts[0]
	}

	last := len(ts) - 1
	l := cons(ts[last-1], ts[last])

	for i := last - 2; i >= 0; i-- {
		l = cons(ts[i], Normal(l))
	}

	return l
})

func cons(t1, t2 *Thunk) listType {
	return listType{t1, t2}
}

var First = NewStrictFunction(func(os ...Object) Object {
	if len(os) != 1 {
		return NumArgsError("first", "1")
	}

	o := os[0]
	l, ok := o.(listType)

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
	l, ok := o.(listType)

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
	l, ok := o.(listType)

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

// ordered

func (l1 listType) less(ord ordered) bool {
	l2 := ord.(listType)

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
