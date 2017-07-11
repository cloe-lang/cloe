package compile

import (
	"github.com/tisp-lang/tisp/src/lib/core"
	"github.com/tisp-lang/tisp/src/lib/scalar"
	"github.com/tisp-lang/tisp/src/lib/std"
)

func prelude() environment {
	e := newEnvironment(scalar.Convert)

	for _, nv := range []struct {
		name  string
		value *core.Thunk
	}{
		{"if", core.If},
		{"$if", core.If},

		{"partial", core.Partial},

		{"first", core.First},
		{"rest", core.Rest},
		{"prepend", core.Prepend},

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
		{"$=", std.Equal},
		{"<", std.Less},
		{"<=", std.LessEq},
		{">", std.Greater},
		{">=", std.GreaterEq},

		{"toStr", core.ToString},
		{"dump", core.Dump},

		{"delete", core.Delete},
		{"$delete", core.Delete},
		{"include", core.Include},
		{"$include", core.Include},
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
