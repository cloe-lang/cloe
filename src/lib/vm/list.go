package vm

type List struct {
	first *Thunk
	rest  *Thunk
}

var emptyList = List{nil, nil}

func NewList(ts ...*Thunk) List {
	return prepend(append(ts, Normal(emptyList))...).(List)
}

func ListThunk(ts ...*Thunk) *Thunk {
	return Normal(NewList(ts...))
}

func (l1 List) Equal(e Equalable) Object {
	l2 := e.(List)

	if l1 == emptyList || l2 == emptyList {
		return l1 == l2
	}

	var bs [2]Bool

	for i, o := range []Object{
		Equal(l1.first, l2.first),
		Equal(l1.rest, l2.rest),
	} {
		b, ok := o.(Bool)

		if !ok {
			return notBoolError(o)
		}

		bs[i] = b
	}

	return bs[0] && bs[1]
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

	return l.first.Eval()
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

	return l.rest.Eval()
}

func notListError(o Object) Error {
	return TypeError(o, "List")
}

func emptyListError() Error {
	return ValueError("The list is empty. You cannot apply rest.")
}
