package desugar

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

func convertCases(v string, cs []ast.Case) (interface{}, []interface{}) {
	var body interface{}
	ls := []interface{}{}

	for _, cs := range groupCases(cs) {
		var nestedLs []interface{}
		body, nestedLs = matchCasesOfSamePatterns(v, cs)
		ls = append(ls, nestedLs...)
	}

	return body, ls
}

func matchCasesOfSamePatterns(v string, cs []ast.Case) (interface{}, []interface{}) {
	// TODO: Implement this function.

	switch getPatternType(cs[0].Pattern()) {
	case listPattern:
	case dictPattern:
	case scalarPattern:
	case namePattern:
	}

	panic("Not implemented")
}

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
