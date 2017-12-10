package compile

import (
	"github.com/tisp-lang/tisp/src/lib/builtins"
	"github.com/tisp-lang/tisp/src/lib/core"
	"github.com/tisp-lang/tisp/src/lib/desugar"
	"github.com/tisp-lang/tisp/src/lib/parse"
	"github.com/tisp-lang/tisp/src/lib/scalar"
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

		{"=", builtins.Equal},
		{"<", builtins.Less},
		{"<=", builtins.LessEq},
		{">", builtins.Greater},
		{">=", builtins.GreaterEq},

		{"toStr", core.ToString},
		{"dump", core.Dump},

		{"delete", core.Delete},
		{"include", core.Include},
		{"insert", core.Insert},
		{"merge", core.Merge},
		{"size", core.Size},
		{"toList", core.ToList},

		{"dict", builtins.Dictionary},
		{"error", core.Error},

		{"y", builtins.Y},
		{"ys", builtins.Ys},

		{"par", builtins.Par},
		{"seq", builtins.Seq},
		{"rally", builtins.Rally},

		{"read", builtins.Read},
		{"write", builtins.Write},

		{"matchError", core.NewError("MatchError", "A value didn't match with any pattern.")},
		{"catch", core.Catch},

		{"pure", core.Pure},
	} {
		e.set(nv.name, nv.value)
		e.set("$"+nv.name, nv.value)
	}

	return e
}()

func builtinsEnvironment() environment {
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
		for n, t := range compileBuiltinModule(e.copy(), "<builtins>", s) {
			e.set(n, t)
			e.set("$"+n, t)
		}
	}

	return e
}

func compileBuiltinModule(e environment, path, source string) module {
	m, err := parse.SubModule(path, source)

	if err != nil {
		panic(err)
	}

	c := newCompiler(e, nil)
	c.compileModule(desugar.Desugar(m))

	return c.env.toMap()
}
