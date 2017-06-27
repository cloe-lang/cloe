package match

import (
	"fmt"

	"github.com/tisp-lang/tisp/src/lib/ast"
	"github.com/tisp-lang/tisp/src/lib/debug"
	"github.com/tisp-lang/tisp/src/lib/gensym"
	"github.com/tisp-lang/tisp/src/lib/scalar"
)

type desugarer struct {
	lets []interface{}
}

func newDesugarer() *desugarer {
	return &desugarer{nil}
}

func (d *desugarer) desugarMatchExpression(x interface{}) interface{} {
	switch x := x.(type) {
	case ast.App:
		return ast.NewApp(
			d.desugarMatchExpression(x.Function()),
			d.desugarMatchExpression(x.Arguments()).(ast.Arguments),
			x.DebugInfo())
	case ast.Arguments:
		ps := make([]ast.PositionalArgument, 0, len(x.Positionals()))

		for _, p := range x.Positionals() {
			ps = append(ps, d.desugarMatchExpression(p).(ast.PositionalArgument))
		}

		ks := make([]ast.KeywordArgument, 0, len(x.Keywords()))

		for _, k := range x.Keywords() {
			ks = append(ks, d.desugarMatchExpression(k).(ast.KeywordArgument))
		}

		dicts := make([]interface{}, 0, len(x.ExpandedDicts()))

		for _, dict := range x.ExpandedDicts() {
			dicts = append(dicts, d.desugarMatchExpression(dict))
		}

		return ast.NewArguments(ps, ks, dicts)
	case ast.KeywordArgument:
		return ast.NewKeywordArgument(x.Name(), d.desugarMatchExpression(x.Value()))
	case ast.LetFunction:
		ls := make([]interface{}, 0, len(x.Lets()))

		for _, l := range x.Lets() {
			l := d.desugarMatchExpression(l)
			ls = append(ls, append(d.lets, l)...)
			d.lets = nil
		}

		return ast.NewLetFunction(
			x.Name(),
			x.Signature(),
			ls,
			d.desugarMatchExpression(x.Body()),
			x.DebugInfo())
	case ast.LetVar:
		return ast.NewLetVar(x.Name(), d.desugarMatchExpression(x.Expr()))
	case ast.Match:
		return app(d.createMatchFunction(x.Cases()), d.desugarMatchExpression(x.Value()))
	case ast.Output:
		return ast.NewOutput(d.desugarMatchExpression(x.Expr()), x.Expanded())
	case ast.PositionalArgument:
		return ast.NewPositionalArgument(d.desugarMatchExpression(x.Value()), x.Expanded())
	default:
		return x
	}
}

func (d *desugarer) createMatchFunction(cs []ast.Case) interface{} {
	arg := gensym.GenSym("match", "argument")
	body := d.casesToBody(arg, cs)

	f := ast.NewLetFunction(
		gensym.GenSym("match", "function"),
		ast.NewSignature([]string{arg}, nil, "", nil, nil, ""),
		d.lets,
		body,
		debug.NewGoInfo(0))

	d.lets = []interface{}{f}

	return f.Name()
}

func (d *desugarer) casesToBody(arg string, cs []ast.Case) interface{} {
	cs = renameBoundNamesInCases(cs)
	body := app("error", "MatchError", "\"Failed to match a value with patterns.\"")

	for _, cs := range groupCases(cs) {
		result, ok := matchCasesOfSamePatterns(arg, cs)

		body = app(
			"if",
			ok,
			result,
			body)
	}

	return body
}

func renameBoundNamesInCases(cs []ast.Case) []ast.Case {
	panic("Not implemented")
}

func app(xs ...interface{}) interface{} {
	return ast.NewPApp(xs[0], xs[1:], debug.NewGoInfo(0))
}

func matchCasesOfSamePatterns(v string, cs []ast.Case) (interface{}, interface{}) {
	// TODO: Implement this function.

	switch getPatternType(cs[0].Pattern()) {
	case listPattern:
		// matchType(v, "list")
		panic("Not implemented")
	case dictPattern:
		panic("Not implemented")
	case scalarPattern:
		panic("Not implemented")
	case namePattern:
		panic("Not implemented")
	}

	panic(fmt.Errorf("Invalid cases: %#v", cs))
}

func matchType(v string, typ string, then interface{}, els interface{}) interface{} {
	return app(
		"if",
		app("=", app("typeOf", v), typ),
		then,
		els)
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
