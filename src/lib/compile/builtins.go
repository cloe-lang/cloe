package compile

import (
	"github.com/tisp-lang/tisp/src/lib/core"
	"github.com/tisp-lang/tisp/src/lib/scalar"
	"github.com/tisp-lang/tisp/src/lib/std"
)

var goBuiltins = func() environment {
	e := newEnvironment(scalar.Convert)

	for _, nv := range []struct {
		name  string
		value *core.Thunk
	}{
		{"if", core.If},

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

		{"dict", std.Dictionary},
		{"error", core.Error},

		{"y", std.Y},
		{"ys", std.Ys},

		{"par", std.Par},
		{"seq", std.Seq},
		{"rally", std.Rally},

		{"read", std.Read},
		{"write", std.Write},

		{"matchError", core.NewError("MatchError", "A value didn't match with any pattern.")},
		{"catch", core.Catch},

		{"pure", core.Pure},
	} {
		e.set(nv.name, nv.value)
		e.set("$"+nv.name, nv.value)
	}

	return e
}()

func builtins() environment {
	e := goBuiltins.copy()

	for _, s := range []string{
		`
			(def (list ..xs) xs)
		`,
		`
			(def (indexOf list elem . (index 0))
				(match list
					[] (error "ElementNotFoundError" "Could not find an element in a list")
					[first ..rest] (if (= first elem) index (indexOf rest elem . index (+ index 1)))))

			(def (map func list)
				(match list
					[] []
					[first ..rest] (prepend (func first) (map func rest))))
		`,
	} {
		for n, t := range subModule(e.copy(), "<builtins>", s) {
			e.set(n, t)
			e.set("$"+n, t)
		}
	}

	return e
}
