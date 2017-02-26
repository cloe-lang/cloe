package ir

import "github.com/raviqqe/tisp/src/lib/core"

type KeywordArgument struct {
	name  string
	value interface{}
}

func NewKeywordArgument(n string, ir interface{}) KeywordArgument {
	return KeywordArgument{
		name:  n,
		value: ir,
	}
}

func (k KeywordArgument) compile(args []*core.Thunk) core.KeywordArgument {
	return core.NewKeywordArgument(k.name, compileExpression(args, k.value))
}
