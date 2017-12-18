package compile

import (
	"github.com/coel-lang/coel/src/lib/builtins"
	"github.com/coel-lang/coel/src/lib/core"
	"github.com/coel-lang/coel/src/lib/desugar"
	"github.com/coel-lang/coel/src/lib/parse"
	"github.com/coel-lang/coel/src/lib/scalar"
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
		{"eseq", builtins.EffectSeq},
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

			(def (generateMaxOrMinFunction name compare)
				(def (maxOrMin ..ns)
					(let x (ns 0))
					(let y (ns 1))
					(let m (if (compare x y) x y))
					(match ns
						[]
							(error
								"ValueError"
								(merge "Number of arguments to " name " function must be greater than 1."))
						[x] x
						[_ _ ..xs] (seq m (maxOrMin m ..xs))))
				maxOrMin)

			(let max (generateMaxOrMinFunction "max" >))
			(let min (generateMaxOrMinFunction "min" <))

			(def (slice list (start 0) (end nil))
				(let end (if (= end nil) (max start (size list)) end))
				(if
					(< end start) (error "ValueError" "start index must be less than end index")
					(= end 0) []
					(match list
						[] []
						[x ..xs]
							(if
								(> start 0) (slice xs (- start 1) (- end 1))
								(= start 0) (prepend x (slice xs 0 (- end 1)))
								[]))))

			(def (generateAndOrOrFunction name operator)
				(def (f ..bools)
					(match bools
						[] (error
							"ValueError"
							(merge "Number of arguments to " name " function must be greater than 1"))
						[x] x
						[x y ..xs] (f (operator x y) ..xs)))
				f)

			(let and (generateAndOrOrFunction "and" (\ (x y) (if x y false))))
			(let or (generateAndOrOrFunction "or" (\ (x y) (if x true y))))

			(def (not bool)
				(if bool false true))

			(def (zip ..lists)
				(if (or ..(map (\ (list) (= 0 (size list))) lists))
					[]
					(prepend (map first lists) (zip ..(map rest lists)))))
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

	_, err = c.compileModule(desugar.Desugar(m))

	if err != nil {
		panic(err)
	}

	return c.env.toMap()
}
