package signature

import ".."

type KeywordArgument struct {
	key   string
	value *vm.Thunk
}

func NewKeywordArgument(key string, value *vm.Thunk) KeywordArgument {
	return KeywordArgument{key: key, value: value}
}
