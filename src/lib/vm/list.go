package vm

type List struct {
	first *Thunk
	rest  *Thunk
}

func NewList(ts ...*Thunk) *Thunk {
	list := Cons(ts[len(ts)-1], Nil())

	for i := len(ts) - 2; i >= 0; i-- {
		list = Cons(ts[i], list)
	}

	return list
}

func Cons(t1, t2 *Thunk) *Thunk {
	return Normal(List{t1, t2})
}

func First(t *Thunk) *Thunk {
	return applyList(func(l List) *Thunk { return l.first }, t)
}

func Rest(t *Thunk) *Thunk {
	return applyList(func(l List) *Thunk { return l.rest }, t)
}

func applyList(f func(List) *Thunk, t *Thunk) *Thunk {
	o := evalList(t)

	if l, ok := o.(List); ok {
		return f(l)
	}

	return Normal(o)
}

func evalList(t *Thunk) Object {
	list, ok := t.Eval().(List)

	if !ok {
		return NewError("Expected List but %#v.", t.Result)
	}

	return list
}
