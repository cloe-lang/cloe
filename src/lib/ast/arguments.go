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

func (a Arguments) Positionals() []PositionalArgument {
	return a.positionals
}

func (a Arguments) Keywords() []KeywordArgument {
	return a.keywords
}

func (a Arguments) ExpandedDicts() []interface{} {
	return a.expandedDicts
}
