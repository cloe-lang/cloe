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
