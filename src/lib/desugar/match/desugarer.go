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
			cs = append(cs, renameBoundNamesInCase(ast.NewMatchCase(c.Pattern(), d.desugar(c.Value()))))
		}

		return app(d.createMatchFunction(cs), d.desugar(x.Value()))
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

func (d *desugarer) createMatchFunction(cs []ast.MatchCase) interface{} {
	arg := gensym.GenSym("match", "argument")
	body := d.desugarCases(arg, cs)

	f := ast.NewLetFunction(
		gensym.GenSym("match", "function"),
		ast.NewSignature([]string{arg}, nil, "", nil, nil, ""),
		d.takeLets(),
		body,
		debug.NewGoInfo(0))

	d.lets = []interface{}{f}

	return f.Name()
}

func (d *desugarer) desugarCases(v interface{}, cs []ast.MatchCase) interface{} {
	css := groupCases(cs)
	ks := []ast.SwitchCase{}

	if cs, ok := css[listPattern]; ok {
		ks = append(ks, ast.NewSwitchCase("list", d.desugarListCases(v, cs)))
	}

	if cs, ok := css[dictPattern]; ok {
		ks = append(ks, ast.NewSwitchCase("dict", d.desugarDictCases(v, cs)))
	}

	dc := interface{}(nil)

	if cs, ok := css[namePattern]; ok {
		dc = d.desugarNameCases(v, cs)
	}

	return ast.NewSwitch(app("typeOf", v), ks, dc)
}

func groupCases(cs []ast.MatchCase) map[patternType][]ast.MatchCase {
	css := map[patternType][]ast.MatchCase{}

	for _, c := range cs {
		t := getPatternType(c.Pattern())
		css[t] = append(css[t], c)
	}

	return css
}

func getPatternType(p interface{}) patternType {
	switch x := p.(type) {
	case string:
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

func (d *desugarer) desugarListCases(v interface{}, cs []ast.MatchCase) interface{} {
	panic("Not implemented")
}

func (d *desugarer) desugarDictCases(v interface{}, cs []ast.MatchCase) interface{} {
	panic("Not implemented")
}

func (d *desugarer) desugarNameCases(v interface{}, cs []ast.MatchCase) interface{} {
	dc := interface{}(nil)
	ks := []ast.SwitchCase{}

	for i, c := range cs {
		p := c.Pattern().(string)
		isScalar := scalar.Defined(p)

		if !isScalar && i < len(cs)-1 {
			panic(fmt.Errorf("A wildcard pattern is found, but some cases are left"))
		} else if !isScalar {
			d.letVar(p, v)
			dc = c.Value()
			break
		}

		ks = append(ks, ast.NewSwitchCase(p, c.Value()))
	}

	return ast.NewSwitch(v, ks, dc)
}

func renameBoundNamesInCase(c ast.MatchCase) ast.MatchCase {
	p, ns := newPatternRenamer().rename(c.Pattern())
	return ast.NewMatchCase(p, newValueRenamer(ns).rename(c.Value()))
}

// func (d *desugarer) casesToBody(arg string, cs []ast.MatchCase) interface{} {
//	cs = renameBoundNamesInCases(cs)
//	body := app("error", "\"MatchError\"", "\"Failed to match a value with patterns.\"")

//	for _, cs := range groupCases(cs) {
//		result, ok := d.matchCasesOfSamePatterns(arg, cs)
//		body = app("if", ok, result, body)
//	}

//	return body
// }

// func (d *desugarer) matchCasesOfSamePatterns(v interface{}, cs []ast.MatchCase) (interface{}, interface{}) {
//	switch getPatternType(cs[0].Pattern()) {
//	case listPattern:
//		panic("Not implemented")
//	case dictPattern:
//		panic("Not implemented")
//	case scalarPattern:
//		ss := make([]interface{}, 0, 2*len(cs))

//		for _, c := range cs {
//			ss = append(ss, c.Pattern(), c.Value())
//		}

//		dict := gensym.GenSym("match", "scalar", "dict")
//		d.letVar(dict, app("dict", ss...))

//		return app(dict, v), app("include", dict, v)
//	case namePattern:
//		if len(cs) != 1 {
//			panic(fmt.Errorf("Duplicate name patterns: %v", len(cs)))
//		}

//		c := cs[0]
//		d.letVar(c.Pattern().(string), v)
//		return c.Value(), "true"
//	}

//	panic(fmt.Errorf("Invalid cases: %#v", cs))
// }

// func matchType(v string, typ string) interface{} {
//	return app("=", app("typeOf", v), typ)
// }
