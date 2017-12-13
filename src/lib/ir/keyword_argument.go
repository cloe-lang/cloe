package ir

import "github.com/coel-lang/coel/src/lib/core"

// KeywordArgument represents a keyword argument passed to a function.
type KeywordArgument struct {
	name  string
	value interface{}
}

// NewKeywordArgument creates a keyword argument from a bound name and its value.
func NewKeywordArgument(n string, v interface{}) KeywordArgument {
	return KeywordArgument{n, v}
}

func (k KeywordArgument) interpret(args []*core.Thunk) core.KeywordArgument {
	return core.NewKeywordArgument(k.name, interpretExpression(args, k.value))
}
