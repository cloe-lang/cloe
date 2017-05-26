package ir

import "github.com/tisp-lang/tisp/src/lib/core"

// KeywordArgument represents a keyword argument passed to a function.
type KeywordArgument struct {
	name  string
	value interface{}
}

// NewKeywordArgument creates a keyword argument from a bound name and its value.
func NewKeywordArgument(n string, ir interface{}) KeywordArgument {
	return KeywordArgument{
		name:  n,
		value: ir,
	}
}

func (k KeywordArgument) compile(args []*core.Thunk) core.KeywordArgument {
	return core.NewKeywordArgument(k.name, compileExpression(args, k.value))
}
