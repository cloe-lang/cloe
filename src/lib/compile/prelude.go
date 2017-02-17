package compile

import (
	"../vm"
	"fmt"
	"strconv"
)

type preludeEnvironment map[string]*vm.Thunk

var prelude preludeEnvironment = map[string]*vm.Thunk{
	"true":  vm.True,
	"false": vm.False,
	"if":    vm.If,

	"partial": vm.Partial,

	"first":   vm.First,
	"rest":    vm.Rest,
	"prepend": vm.Prepend,

	"nil": vm.Nil,

	"+":   vm.Add,
	"-":   vm.Sub,
	"*":   vm.Mul,
	"/":   vm.Div,
	"mod": vm.Mod,
	"pow": vm.Pow,

	"y":  vm.Y,
	"ys": vm.Ys,

	"cause": vm.Cause,

	"write": write,
}

func (e preludeEnvironment) get(s string) *vm.Thunk {
	t, ok := e[s]

	if ok {
		return t
	}

	f, err := strconv.ParseFloat(s, 64)

	if err == nil {
		return vm.NewNumber(f)
	}

	panic(fmt.Sprint(s, "is not found."))
}
