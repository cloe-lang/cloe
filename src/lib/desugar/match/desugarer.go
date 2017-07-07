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

func (d *desugarer) desugar(x interface{}) interface{} {
	switch x := x.(type) {
	case ast.App:
		return ast.NewApp(
			d.desugar(x.Function()),
			d.desugar(x.Arguments()).(ast.Arguments),
			x.DebugInfo())
	case ast.Arguments:
		ps := make([]ast.PositionalArgument, 0, len(x.Positionals()))

		for _, p := range x.Positionals() {
			ps = append(ps, d.desugar(p).(ast.PositionalArgument))
		}

		ks := make([]ast.KeywordArgument, 0, len(x.Keywords()))

		for _, k := range x.Keywords() {
			ks = append(ks, d.desugar(k).(ast.KeywordArgument))
		}

		dicts := make([]interface{}, 0, len(x.ExpandedDicts()))

		for _, dict := range x.ExpandedDicts() {
			dicts = append(dicts, d.desugar(dict))
		}

		return ast.NewArguments(ps, ks, dicts)
	case ast.KeywordArgument:
		return ast.NewKeywordArgument(x.Name(), d.desugar(x.Value()))
	case ast.LetFunction:
		ls := make([]interface{}, 0, len(x.Lets()))

		for _, l := range x.Lets() {
			l := d.desugar(l)
			ls = append(ls, append(d.takeLets(), l)...)
		}

		b := d.desugar(x.Body())
		ls = append(ls, d.takeLets()...)

		return ast.NewLetFunction(
			x.Name(),
			x.Signature(),
			ls,
			b,
			x.DebugInfo())
	case ast.LetVar:
		return ast.NewLetVar(x.Name(), d.desugar(x.Expr()))
	case ast.Match:
		cs := make([]ast.MatchCase, 0, len(x.Cases()))

		for _, c := range x.Cases() {
			cs = append(cs, ast.NewMatchCase(c.Pattern(), d.desugar(c.Value())))
		}

		return d.desugarMatchExpression(ast.NewMatch(d.desugar(x.Value()), cs))
	case ast.Output:
		return ast.NewOutput(d.desugar(x.Expr()), x.Expanded())
	case ast.PositionalArgument:
		return ast.NewPositionalArgument(d.desugar(x.Value()), x.Expanded())
	default:
		return x
	}
}

func (d *desugarer) takeLets() []interface{} {
	ls := d.lets
	d.lets = nil
	return ls
}

func (d *desugarer) letVar(s string, v interface{}) {
	d.lets = append(d.lets, ast.NewLetVar(s, v))
}

func (d *desugarer) desugarMatchExpression(m ast.Match) interface{} {
	panic("Not implemented")
}

func (d *desugarer) casesToBody(arg string, cs []ast.MatchCase) interface{} {
	cs = renameBoundNamesInCases(cs)
	body := app("error", "\"MatchError\"", "\"Failed to match a value with patterns.\"")

	for _, cs := range groupCases(cs) {
		result, ok := d.matchCasesOfSamePatterns(arg, cs)
		body = app("if", ok, result, body)
	}

	return body
}

func renameBoundNamesInCases(cs []ast.MatchCase) []ast.MatchCase {
	new := make([]ast.MatchCase, 0, len(cs))

	for _, c := range cs {
		new = append(new, renameBoundNamesInCase(c))
	}

	return new
}

func renameBoundNamesInCase(c ast.MatchCase) ast.MatchCase {
	p, ns := newPatternRenamer().rename(c.Pattern())
	return ast.NewMatchCase(p, newValueRenamer(ns).rename(c.Value()))
}

func app(f interface{}, args ...interface{}) interface{} {
	return ast.NewPApp(f, args, debug.NewGoInfo(0))
}

func (d *desugarer) matchCasesOfSamePatterns(v interface{}, cs []ast.MatchCase) (interface{}, interface{}) {
	switch getPatternType(cs[0].Pattern()) {
	case listPattern:
		panic("Not implemented")
	case dictPattern:
		panic("Not implemented")
	case scalarPattern:
		ss := make([]interface{}, 0, 2*len(cs))

		for _, c := range cs {
			ss = append(ss, c.Pattern(), c.Value())
		}

		dict := gensym.GenSym("match", "scalar", "dict")
		d.letVar(dict, app("dict", ss...))

		return app(dict, v), app("include", dict, v)
	case namePattern:
		if len(cs) != 1 {
			panic(fmt.Errorf("Duplicate name patterns: %v", len(cs)))
		}

		c := cs[0]
		d.letVar(c.Pattern().(string), v)
		return c.Value(), "true"
	}

	panic(fmt.Errorf("Invalid cases: %#v", cs))
}

// func matchType(v string, typ string) interface{} {
//	return app("=", app("typeOf", v), typ)
// }

func groupCases(cs []ast.MatchCase) map[patternType][]ast.MatchCase {
	m := map[patternType][]ast.MatchCase{}

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
		if len(x.Arguments().Positionals()) == 0 {
			return scalarPattern
		}

		switch x.Function().(string) {
		case "$list":
			return listPattern
		case "$dict":
			return dictPattern
		}
	}

	panic(fmt.Errorf("Invalid pattern: %#v", p))
}
