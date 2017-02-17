package ir

import "../vm"

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

func (k KeywordArgument) compile(args []*vm.Thunk) vm.KeywordArgument {
	return vm.NewKeywordArgument(k.name, compileExpression(args, k.value))
}
