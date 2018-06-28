package compile

import (
	"github.com/cloe-lang/cloe/src/lib/builtins"
	"github.com/cloe-lang/cloe/src/lib/core"
	"github.com/cloe-lang/cloe/src/lib/desugar"
	"github.com/cloe-lang/cloe/src/lib/parse"
	"github.com/cloe-lang/cloe/src/lib/scalar"
)

const builtinsFilename = "<builtins>"

var goBuiltins = func() environment {
	e := newEnvironment(scalar.Convert)

	for s, t := range map[string]core.Value{
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

		"=":  core.Equal,
		"<":  builtins.Less,
		"<=": builtins.LessEq,
		">":  builtins.Greater,
		">=": builtins.GreaterEq,

		"toString": core.ToString,
		"dump":     core.Dump,

		"@":       core.Index,
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
		"print": builtins.Print,

		"error": core.Error,
		"catch": core.Catch,

		"pure": core.Pure,
	} {
		e.set(s, t)
		e.set("$"+s, t)
	}

	for s, t := range map[string]core.Value{
		"matchError": core.NewError("MatchError", "a value didn't match with any pattern"),
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

			(def (dictionary ..args)
				(if (= args [])
					{}
					(insert
						(dictionary ..(rest (rest args)))
						(first args)
						(first (rest args)))))
		`) {
		e.set("$"+n, t)
	}

	for n, t := range compileBuiltinModule(e.copy(), builtinsFilename, `
			(def (boolean? x) (= (typeOf x) "boolean"))
			(def (dictionary? x) (= (typeOf x) "dictionary"))
			(def (function? x) (= (typeOf x) "function"))
			(def (list? x) (= (typeOf x) "list"))
			(def (nil? x) (= (typeOf x) "nil"))
			(def (number? x) (= (typeOf x) "number"))
			(def (string? x) (= (typeOf x) "string"))

			(def (index list elem . i 1)
				(match list
					[] (error "ElementNotFoundError" "Could not find an element in a list")
					[first ..rest] (if (= first elem) i (index rest elem . i (+ i 1)))))

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
							(let m (if (compare x y) x y)
								(seq m (maxOrMin m ..xs)))))
				maxOrMin)

			(let max (generateMaxOrMinFunction "max" >))
			(let min (generateMaxOrMinFunction "min" <))

			(def (slice list . start 1 end nil)
				(let end (if (= end nil) (max start (size list)) end))
				(if
					(string? list) (merge "" ..(slice (toList list) . start start end end))
					(< end start) (error "ValueError" "start index must be less than end index")
					(match list
						[] []
						[x ..xs]
							(if
								(= end 1) [x]
								(> start 1) (slice xs . start (- start 1) end (- end 1))
								(= start 1) [x ..(slice xs . start 1 end (- end 1))]
								[]))))

			(def (generateAndOrOrFunction name operator)
				(def (f ..bs)
					(match bs
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

			(def (filter func list)
				(match list
					[] []
					[x ..xs]
						(match (filter func xs)
							xs (if (func x) [x ..xs] xs))))

			(def (sort list . less <)
				(let sort (partial sort . less less))
				(let pivot (@ list (// (size list) 2)))
				(match list
					[] []
					[x] [x]
					_ [
						..(sort (filter (\ (x) (less x pivot)) list))
						..(filter (\ (x) (and (not (less x pivot)) (not (less pivot x)))) list)
						..(sort (filter (\ (x) (less pivot x)) list))]))

			(let sortImpl sort)

			(def (sort list . less <)
				(let x (first list))
				(match list
					[] []
					_ (if (less x x)
						(error "ValueError" "less function should not be reflexive.")
						(sortImpl list . less less))))
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

	_, err = c.compileModule(desugar.Desugar(m), "INVALID")

	if err != nil {
		panic(err)
	}

	return c.env.toMap()
}
