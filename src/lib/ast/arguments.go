package ast

// Arguments represents arguments passed to a function.
type Arguments struct {
	positionals   []PositionalArgument
	keywords      []KeywordArgument
	expandedDicts []interface{}
}

// NewArguments creates arguments.
func NewArguments(
	ps []PositionalArgument,
	ks []KeywordArgument,
	expandedDicts []interface{}) Arguments {
	return Arguments{ps, ks, expandedDicts}
}

// Positionals returns positional arguments contained in arguments.
func (a Arguments) Positionals() []PositionalArgument {
	return a.positionals
}

// Keywords returns keyword arguments contained in arguments.
func (a Arguments) Keywords() []KeywordArgument {
	return a.keywords
}

// ExpandedDicts returns expanded dictionary arguments contained in arguments.
func (a Arguments) ExpandedDicts() []interface{} {
	return a.expandedDicts
}
