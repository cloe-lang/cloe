package ir

// Arguments represents arguments passed to a function in IR.
type Arguments struct {
	positionals   []PositionalArgument
	keywords      []KeywordArgument
	expandedDicts []interface{}
}

// NewArguments creates arguments from positional and keyword arguments and
// expanded dictionaries.
func NewArguments(ps []PositionalArgument, ks []KeywordArgument, dicts []interface{}) Arguments {
	return Arguments{
		positionals:   ps,
		keywords:      ks,
		expandedDicts: dicts,
	}
}
