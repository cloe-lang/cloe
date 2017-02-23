package ast

type KeywordArgument struct {
	name  string
	value interface{}
}

func NewKeywordArgument(name string, value interface{}) KeywordArgument {
	return KeywordArgument{name, value}
}

func (k KeywordArgument) Name() string {
	return k.name
}

func (k KeywordArgument) Value() interface{} {
	return k.value
}
