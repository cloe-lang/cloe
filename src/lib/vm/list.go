package vm

type List struct {
	first *Thunk
	rest  *Thunk
}

var rawEmptyList = List{nil, nil}
var EmptyList = Normal(rawEmptyList)

func NewList(ts ...*Thunk) *Thunk {
	return Prepend(append(ts, EmptyList)...)
}

func cons(t1, t2 *Thunk) *Thunk {
	return Normal(List{t1, t2})
}

func (l1 List) Equal(o Object) bool {
	l2 := o.(List)

	if l1 == rawEmptyList || l2 == rawEmptyList {
		return l1 == l2
	}

	return bool(And(
		Equal(l1.first, l2.first),
		Equal(l1.rest, l2.rest)).Eval().(Bool))
}

func Prepend(ts ...*Thunk) *Thunk {
	if len(ts) == 0 {
		return NewError("Too few argument for prepend.")
	}

	last := len(ts) - 1
	l := ts[last]

	for i := last - 1; i >= 0; i-- {
		l = cons(ts[i], l)
	}

	return l
}

func First(ts ...*Thunk) *Thunk {
	if len(ts) != 1 {
		return NumArgsError("first", "1")
	}

	return ts[0].Eval().(List).first
}

func Rest(ts ...*Thunk) *Thunk {
	if len(ts) != 1 {
		return NumArgsError("rest", "1")
	}

	l := ts[0].Eval().(List)

	if l == rawEmptyList {
		return NewError("The list is empty. You cannot apply rest.")
	}

	return l.rest
}
