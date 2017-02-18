package compile

import (
	"../vm"
	"./env"
	"strconv"
)

var prelude = func() env.Environment {
	e := env.NewEnvironment(func(s string) (*vm.Thunk, error) {
		n, err := strconv.ParseFloat(s, 64)

		if err != nil {
			return nil, err
		}

		return vm.NewNumber(n), nil
	})

	for _, nv := range []struct {
		name  string
		value *vm.Thunk
	}{
		{"true", vm.True},
		{"false", vm.False},
		{"if", vm.If},

		{"partial", vm.Partial},

		{"first", vm.First},
		{"rest", vm.Rest},
		{"prepend", vm.Prepend},

		{"nil", vm.Nil},

		{"+", vm.Add},
		{"-", vm.Sub},
		{"*", vm.Mul},
		{"/", vm.Div},
		{"mod", vm.Mod},
		{"pow", vm.Pow},

		{"y", vm.Y},
		{"ys", vm.Ys},

		{"cause", vm.Cause},

		{"write", write},
	} {
		e.Set(nv.name, nv.value)
	}

	return e
}()
