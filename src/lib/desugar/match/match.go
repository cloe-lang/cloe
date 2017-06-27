package match

import (
	"fmt"

	"github.com/tisp-lang/tisp/src/lib/ast"
	"github.com/tisp-lang/tisp/src/lib/debug"
	"github.com/tisp-lang/tisp/src/lib/gensym"
	"github.com/tisp-lang/tisp/src/lib/scalar"
)

type patternType int

const (
	listPattern patternType = iota
	dictPattern
	scalarPattern
	namePattern
)

// Desugar desugars match expressions in an AST node.
func Desugar(x interface{}) []interface{} {
	return []interface{}{desugarMatchExpression(x)}
}

func desugarMatchExpression(x interface{}) interface{} {
	switch x := x.(type) {
	case ast.App:
		return ast.NewApp(
			desugarMatchExpression(x.Function()),
			desugarMatchExpression(x.Arguments()).(ast.Arguments),
			x.DebugInfo())
	case ast.Arguments:
		ps := make([]ast.PositionalArgument, 0, len(x.Positionals()))

		for _, p := range x.Positionals() {
			ps = append(ps, desugarMatchExpression(p).(ast.PositionalArgument))
		}

		ks := make([]ast.KeywordArgument, 0, len(x.Keywords()))

		for _, k := range x.Keywords() {
			ks = append(ks, desugarMatchExpression(k).(ast.KeywordArgument))
		}

		ds := make([]interface{}, 0, len(x.ExpandedDicts()))

		for _, d := range x.ExpandedDicts() {
			ds = append(ds, desugarMatchExpression(d))
		}

		return ast.NewArguments(ps, ks, ds)
	case ast.KeywordArgument:
		return ast.NewKeywordArgument(x.Name(), desugarMatchExpression(x.Value()))
	case ast.LetFunction:
		ls := make([]interface{}, 0, len(x.Lets()))

		for _, l := range x.Lets() {
			ls = append(ls, desugarMatchExpression(l))
		}

		return ast.NewLetFunction(
			x.Name(),
			x.Signature(),
			ls,
			desugarMatchExpression(x.Body()),
			x.DebugInfo())
	case ast.LetVar:
		return ast.NewLetVar(x.Name(), desugarMatchExpression(x.Expr()))
	case ast.Match:
		return desugarMatchIntoApp(x)
	case ast.Output:
		return ast.NewOutput(desugarMatchExpression(x.Expr()), x.Expanded())
	case ast.PositionalArgument:
		return ast.NewPositionalArgument(desugarMatchExpression(x.Value()), x.Expanded())
	default:
		return x
	}
}

func desugarMatchIntoApp(m ast.Match) interface{} {
	return ast.NewApp(
		createMatchFunction(m.Cases()),
		ast.NewArguments(
			[]ast.PositionalArgument{ast.NewPositionalArgument(m.Value(), false)},
			nil,
			nil),
		debug.NewGoInfo(0))
}

func createMatchFunction(cs []ast.Case) interface{} {
	arg := gensym.GenSym("match", "argument")
	body, ls := convertCases(arg, cs)

	return ast.NewLetFunction(
		gensym.GenSym("match", "function"),
		ast.NewSignature([]string{arg}, nil, "", nil, nil, ""),
		ls,
		body,
		debug.NewGoInfo(0))
}

func convertCases(matched interface{}, cs []ast.Case) (interface{}, []interface{}) {
	fs := []interface{}{}

	for _, cs := range groupCases(cs) {
		fs = append(fs, createFunctionOfSameCases(cs))
	}

	ls := []interface{}{}

	body := app("error", "MatchError", "\"Failed to match a value with patterns.\"")
	for _, f := range fs {
		result := gensym.GenSym("match", "intermediate", "result")
		value := gensym.GenSym("match", "intermediate", "value")
		ok := gensym.GenSym("match", "intermediate", "ok")

		ls = append(
			ls,
			ast.NewLetVar(result, app(f, matched)),
			ast.NewLetVar(value, app("first", result)),
			ast.NewLetVar(ok, app("first", app("rest", result))))

		body = app("if", ok, value, body)
	}

	return body, append(fs, ls...)
}

func app(xs ...interface{}) interface{} {
	return ast.NewPApp(xs[0], xs[1:], debug.NewGoInfo(0))
}

func createFunctionOfSameCases(cs []ast.Case) interface{} {
	// Functions should return [result, ok: bool].
	panic("Not implemented")
}

// func matchCasesOfSamePatterns(v string, cs []ast.Case) (interface{}, []interface{}) {
//	// TODO: Implement this function.

//	switch getPatternType(cs[0].Pattern()) {
//	case listPattern:
//		matchType(v, "list")
//	case dictPattern:
//	case scalarPattern:
//	case namePattern:
//		return nil, nil
//	}

//	panic("Not implemented")
// }

// func matchType(v string, typ string, then interface{}, els interface{}) interface{} {
//	return ast.NewPApp(
//		"if",
//		[]interface{}{
//			ast.NewPApp("=", []interface{}{ast.NewPApp("typeOf", v), typ}),
//			then,
//			els},
//		debug.NewGoInfo(0))
// }

func groupCases(cs []ast.Case) map[patternType][]ast.Case {
	m := map[patternType][]ast.Case{}

	for _, c := range cs {
		p := getPatternType(c.Pattern())
		m[p] = append(m[p], c)
	}

	return m
}

func getPatternType(p interface{}) patternType {
	switch x := p.(type) {
	case string:
		if scalar.Defined(x) {
			return scalarPattern
		}

		return namePattern
	case ast.App:
		switch x.Function().(string) {
		case "$list":
			return listPattern
		case "$dict":
			return dictPattern
		}
	}

	panic(fmt.Errorf("Invalid pattern: %#v", p))
}
