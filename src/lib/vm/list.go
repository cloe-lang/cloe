package vm

type List struct {
	first *Thunk
	rest  *Thunk
}

var emptyList = List{nil, nil}

func NewList(ts ...*Thunk) List {
	return App(Normal(Prepend), append(ts, Normal(emptyList))...).Eval().(List)
}

func ListThunk(ts ...*Thunk) *Thunk {
	return Normal(NewList(ts...))
}

func (l1 List) Equal(e Equalable) Object {
	l2 := e.(List)

	if l1 == emptyList || l2 == emptyList {
		return rawBool(l1 == l2)
	}

	for _, t := range []*Thunk{
		// Don't evaluate these parallely for short circuit behavior.
		App(Normal(Equal), l1.first, l2.first),
		App(Normal(Equal), l1.rest, l2.rest),
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

var Prepend = NewLazyFunction(prepend)

func prepend(ts ...*Thunk) Object {
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
}

func cons(t1, t2 *Thunk) List {
	return List{t1, t2}
}

var First = NewStrictFunction(first)

func first(os ...Object) Object {
	if len(os) != 1 {
		return NumArgsError("first", "1")
	}

	o := os[0]
	l, ok := o.(List)

	if !ok {
		return notListError(o)
	} else if l == emptyList {
		return emptyListError()
	}

	return l.first
}

var Rest = NewStrictFunction(rest)

func rest(os ...Object) Object {
	if len(os) != 1 {
		return NumArgsError("rest", "1")
	}

	o := os[0]
	l, ok := o.(List)

	if !ok {
		return notListError(o)
	} else if l == emptyList {
		return emptyListError()
	}

	return l.rest
}

func notListError(o Object) *Thunk {
	return TypeError(o, "List")
}

func emptyListError() *Thunk {
	return ValueError("The list is empty. You cannot apply rest.")
}
