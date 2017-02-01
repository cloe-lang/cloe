package signature

import "../vm"

type Arguments struct {
	positionals   []PositionalArgument
	keywords      []KeywordArgument
	expandedDicts []*vm.Thunk
}

func NewArguments(
	ps []PositionalArgument,
	ks []KeywordArgument,
	expandedDicts []*vm.Thunk) Arguments {
	return Arguments{
		positionals:   ps,
		keywords:      ks,
		expandedDicts: expandedDicts,
	}
}

type PositionalArgument struct {
	value    *vm.Thunk
	expanded bool
}

func NewPositionalArgument(value *vm.Thunk, expanded bool) PositionalArgument {
	return PositionalArgument{
		value:    value,
		expanded: expanded,
	}
}

type KeywordArgument struct {
	key   string
	value *vm.Thunk
}

func NewKeywordArgument(key string, value *vm.Thunk) KeywordArgument {
	return KeywordArgument{key: key, value: value}
}
