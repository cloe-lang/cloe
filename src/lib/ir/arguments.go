package ir

type Arguments struct {
	positionals   []PositionalArgument
	keywords      []KeywordArgument
	expandedDicts []interface{}
}

func NewArguments(ps []PositionalArgument, ks []KeywordArgument, dicts []interface{}) Arguments {
	return Arguments{
		positionals:   ps,
		keywords:      ks,
		expandedDicts: dicts,
	}
}

func NewPositionalArguments(xs ...interface{}) Arguments {
	ps := make([]PositionalArgument, len(xs))

	for i, x := range xs {
		ps[i] = NewPositionalArgument(x, false)
	}

	return NewArguments(ps, []KeywordArgument{}, []interface{}{})
}
