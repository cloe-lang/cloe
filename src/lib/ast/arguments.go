package ast

type Arguments struct {
	positionals   []PositionalArgument
	keywords      []KeywordArgument
	expandedDicts []interface{}
}

func NewArguments(
	ps []PositionalArgument,
	ks []KeywordArgument,
	expandedDicts []interface{}) Arguments {
	return Arguments{ps, ks, expandedDicts}
}
