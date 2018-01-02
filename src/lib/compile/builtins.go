package compile

import (
	"github.com/coel-lang/coel/src/lib/builtins"
	"github.com/coel-lang/coel/src/lib/core"
	"github.com/coel-lang/coel/src/lib/desugar"
	"github.com/coel-lang/coel/src/lib/parse"
	"github.com/coel-lang/coel/src/lib/scalar"
)

const builtinsFilename = "<builtins>"

var goBuiltins = func() environment {
	e := newEnvironment(scalar.Convert)

	for s, t := range map[string]*core.Thunk{
		"if": core.If,

		"partial": core.Partial,

		"first": core.First,
		"rest":  core.Rest,

		"typeOf":   core.TypeOf,
		"ordered?": core.IsOrdered,

		"+":   core.Add,
		"-":   core.Sub,
		"*":   core.Mul,
		"/":   core.Div,
		"//":  core.FloorDiv,
		"mod": core.Mod,
		"**":  core.Pow,

		"=":  builtins.Equal,
		"<":  builtins.Less,
		"<=": builtins.LessEq,
		">":  builtins.Greater,
		">=": builtins.GreaterEq,

		"toString": core.ToString,
		"dump":     core.Dump,

		"delete":  core.Delete,
		"include": core.Include,
		"insert":  core.Insert,
		"merge":   core.Merge,
		"size":    core.Size,
		"toList":  core.ToList,

		"par":   builtins.Par,
		"seq":   builtins.Seq,
		"seq!":  builtins.EffectSeq,
		"rally": builtins.Rally,

		"read":  builtins.Read,
		"write": builtins.Write,

		"error": core.Error,
		"catch": core.Catch,

		"pure": core.Pure,
	} {
		e.set(s, t)
		e.set("$"+s, t)
	}

	for s, t := range map[string]*core.Thunk{
		"matchError": core.NewError("MatchError", "A value didn't match with any pattern."),
		"y":          builtins.Y,
		"ys":         builtins.Ys,
	} {
		e.set("$"+s, t)
	}

	return e
}()

func builtinsEnvironment() environment {
	e := goBuiltins.copy()

	for n, t := range compileBuiltinModule(e.copy(), builtinsFilename, `
			(def (list ..args)
				args)

			(def (dict ..args)
				(if (= args [])
					{}
					(insert
						(dict ..(rest (rest args)))
						(first args)
						(first (rest args)))))
		`) {
		e.set("$"+n, t)
	}

	for n, t := range compileBuiltinModule(e.copy(), builtinsFilename, `
			(def (bool? x) (= (typeOf x) "bool"))
			(def (dict? x) (= (typeOf x) "dict"))
			(def (function? x) (= (typeOf x) "function"))
			(def (list? x) (= (typeOf x) "list"))
			(def (nil? x) (= (typeOf x) "nil"))
			(def (number? x) (= (typeOf x) "number"))
			(def (string? x) (= (typeOf x) "string"))

			(def (indexOf list elem . (index 1))
				(match list
					[] (error "ElementNotFoundError" "Could not find an element in a list")
					[first ..rest] (if (= first elem) index (indexOf rest elem . index (+ index 1)))))

			(def (map func list)
				(match list
					[] []
					[first ..rest] [(func first) ..(map func rest)]))

			(def (reduce func list)
				(match list
					[x] x
					[x y ..xs] (reduce func [(func x y) ..xs])))

			(def (generateMaxOrMinFunction name compare)
				(def (maxOrMin ..args)
					(match args
						[]
							(error
								"ValueError"
								(merge "Number of arguments to " name " function must be greater than 1."))
						[x] x
						[x y ..xs]
							(match (if (compare x y) x y)
								m (seq m (maxOrMin m ..xs)))))
				maxOrMin)

			(let max (generateMaxOrMinFunction "max" >))
			(let min (generateMaxOrMinFunction "min" <))

			(def (slice list (start 1) (end nil))
				(let end (if (= end nil) (max start (size list)) end))
				(if
					(string? list) (merge "" ..(slice (toList list) start end))
					(< end start) (error "ValueError" "start index must be less than end index")
					(match list
						[] []
						[x ..xs]
							(if
								(= end 1) [x]
								(> start 1) (slice xs (- start 1) (- end 1))
								(= start 1) [x ..(slice xs 1 (- end 1))]
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
				(if (or ..(map (\ (list) (= list [])) lists))
					[]
					[(map first lists) ..(zip ..(map rest lists))]))
		`) {
		e.set(n, t)
		e.set("$"+n, t)
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
