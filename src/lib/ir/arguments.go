package ir

// Arguments represents arguments passed to a function in IR.
type Arguments struct {
	positionals []PositionalArgument
	keywords    []KeywordArgument
}

// NewArguments creates arguments from positional and keyword arguments and
// expanded dictionaries.
func NewArguments(ps []PositionalArgument, ks []KeywordArgument) Arguments {
	return Arguments{ps, ks}
}
