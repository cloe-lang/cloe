package parse

import "../types"

func Parse(source string) types.Object {
	o, err := NewState(source).Module()()

	if err != nil {
		panic(err.Error())
	}

	return o
}
