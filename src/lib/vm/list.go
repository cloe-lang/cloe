package vm

type List struct {
	first *Thunk
	rest  *Thunk
}

var EmptyList = Normal(List{nil, nil})

func NewList(ts ...*Thunk) *Thunk {
	l := cons(ts[len(ts)-1], EmptyList)

	for i := len(ts) - 2; i >= 0; i-- {
		l = cons(ts[i], l)
	}

	return l
}

func cons(t1, t2 *Thunk) *Thunk {
	return Normal(List{t1, t2})
}

func Prepend(args List) *Thunk {
	return Normal(List{args.first, args.rest.Eval().(List).first})
}

func First(args List) *Thunk {
	return args.first.Eval().(List).first
}

func Rest(args List) *Thunk {
	l := args.first.Eval().(List)

	if l.first == nil {
		return NewError("The list is empty. You cannot apply rest.")
	} else if l.rest == nil {
		return NewError("The list's length is 1. You cannot apply rest.")
	}

	return l.rest
}
