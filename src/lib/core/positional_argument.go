package core

type PositionalArgument struct {
	value    *Thunk
	expanded bool
}

func NewPositionalArgument(value *Thunk, expanded bool) PositionalArgument {
	return PositionalArgument{
		value:    value,
		expanded: expanded,
	}
}

func NewPositionalArguments(ts ...*Thunk) Arguments {
	ps := make([]PositionalArgument, len(ts))

	for i, t := range ts {
		ps[i] = NewPositionalArgument(t, false)
	}

	return NewArguments(ps, nil, nil)
}
