package core

type KeywordArgument struct {
	name  string
	value *Thunk
}

func NewKeywordArgument(name string, value *Thunk) KeywordArgument {
	return KeywordArgument{name: name, value: value}
}
