package signature

import ".."

type KeywordArgument struct {
	name  string
	value *vm.Thunk
}

func NewKeywordArgument(name string, value *vm.Thunk) KeywordArgument {
	return KeywordArgument{name: name, value: value}
}
