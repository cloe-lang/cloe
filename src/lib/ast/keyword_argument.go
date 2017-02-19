package ast

type KeywordArgument struct {
	name  string
	value interface{}
}

func NewKeywordArgument(name string, value interface{}) KeywordArgument {
	return KeywordArgument{name, value}
}
