package compile

import (
	"fmt"
	"strconv"

	"github.com/tisp-lang/tisp/src/lib/core"
	"github.com/tisp-lang/tisp/src/lib/std"
)

func prelude() environment {
	e := newEnvironment(func(name string) (*core.Thunk, error) {
		if n, err := strconv.ParseInt(name, 0, 64); err == nil && name[0] == '0' {
			return core.NewNumber(float64(n)), nil
		}

		if n, err := strconv.ParseFloat(name, 64); err == nil {
			return core.NewNumber(n), nil
		}

		if s, err := strconv.Unquote(name); err == nil {
			return core.NewString(s), nil
		}

		return nil, fmt.Errorf("the name, %s not found", name)
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

		{"typeOf", core.TypeOf},
		{"isOrdered", core.IsOrdered},

		{"+", core.Add},
		{"-", core.Sub},
		{"*", core.Mul},
		{"/", core.Div},
		{"//", core.FloorDiv},
		{"mod", core.Mod},
		{"**", core.Pow},

		{"=", std.Equal},
		{"<", std.Less},
		{"<=", std.LessEq},
		{">", std.Greater},
		{">=", std.GreaterEq},

		{"toStr", core.ToString},
		{"dump", core.Dump},

		{"delete", core.Delete},
		{"include", core.Include},
		{"insert", core.Insert},
		{"merge", core.Merge},
		{"size", core.Size},
		{"toList", core.ToList},

		{"list", std.List},
		{"$list", std.List},
		{"dict", std.Dictionary},
		{"$dict", std.Dictionary},
		{"error", core.Error},

		{"y", std.Y},
		{"$y", std.Y},
		{"ys", std.Ys},
		{"$ys", std.Ys},

		{"par", std.Par},
		{"seq", std.Seq},
		{"rally", std.Rally},

		{"read", std.Read},
		{"write", std.Write},
	} {
		e.set(nv.name, nv.value)
	}

	return e
}
