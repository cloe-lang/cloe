package core

// KeywordArgument represents a keyword argument passed to a function.
type KeywordArgument struct {
	name  string
	value *Thunk
}

// NewKeywordArgument creates a KeywordArgument from a bound name and its value.
func NewKeywordArgument(name string, value *Thunk) KeywordArgument {
	return KeywordArgument{name: name, value: value}
}
