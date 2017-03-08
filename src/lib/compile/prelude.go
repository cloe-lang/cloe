package compile

import (
	"github.com/raviqqe/tisp/src/lib/compile/env"
	"github.com/raviqqe/tisp/src/lib/core"
	"github.com/raviqqe/tisp/src/lib/std"
	"strconv"
)

var prelude = func() env.Environment {
	e := env.NewEnvironment(func(s string) (*core.Thunk, error) {
		n, err := strconv.ParseFloat(s, 64)

		if err != nil {
			return nil, err
		}

		return core.NewNumber(n), nil
	})

	for _, nv := range []struct {
		name  string
		value *core.Thunk
	}{
		{"true", core.True},
		{"false", core.False},
		{"if", core.If},

		{"partial", core.Partial},

		{"first", core.First},
		{"rest", core.Rest},
		{"prepend", core.Prepend},

		{"nil", core.Nil},

		{"+", core.Add},
		{"-", core.Sub},
		{"*", core.Mul},
		{"/", core.Div},
		{"mod", core.Mod},
		{"pow", core.Pow},

		{"=", core.Equal},
		{"toList", core.ToList},
		{"toStr", core.ToString},

		{"list", std.List},
		{"$list", std.List},
		{"dict", std.Dictionary},
		{"$dict", std.Dictionary},

		{"y", std.Y},
		{"ys", std.Ys},

		{"cause", std.Cause},

		{"write", std.Write},
	} {
		e.Set(nv.name, nv.value)
	}

	return e
}()
