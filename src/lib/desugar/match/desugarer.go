package match

import (
	"fmt"

	"github.com/tisp-lang/tisp/src/lib/ast"
	"github.com/tisp-lang/tisp/src/lib/debug"
	"github.com/tisp-lang/tisp/src/lib/gensym"
	"github.com/tisp-lang/tisp/src/lib/scalar"
)

type desugarer struct {
	letBoundNames, lets []interface{}
}

func newDesugarer() *desugarer {
	return &desugarer{nil, nil}
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
	ls := append(d.letBoundNames, d.lets...)
	d.letBoundNames = nil
	d.lets = nil
	return ls
}

func (d *desugarer) letVar(s string, v interface{}) string {
	d.lets = append(d.lets, ast.NewLetVar(s, v))
	return s
}

func (d *desugarer) bindName(s string, v interface{}) {
	d.letBoundNames = append(d.letBoundNames, ast.NewLetVar(s, v))
}

func (d *desugarer) app(f interface{}, args ...interface{}) string {
	return d.letVar(gensym.GenSym("match", "app"), app(f, args...))
}

func (d *desugarer) createMatchFunction(cs []ast.MatchCase) interface{} {
	arg := gensym.GenSym("match", "argument")
	body := d.desugarCases(arg, cs, "$matchError")

	f := ast.NewLetFunction(
		gensym.GenSym("match", "function"),
		ast.NewSignature([]string{arg}, nil, "", nil, nil, ""),
		d.takeLets(),
		body,
		debug.NewGoInfo(0))

	d.lets = append(d.lets, f)

	return f.Name()
}

func (d *desugarer) desugarCases(v interface{}, cs []ast.MatchCase, dc interface{}) interface{} {
	css := groupCases(cs)

	if cs, ok := css[namePattern]; ok {
		c := cs[0]
		d.bindName(c.Pattern().(string), v)
		dc = c.Value()
	}

	ks := []ast.SwitchCase{}

	if cs, ok := css[listPattern]; ok {
		ks = append(ks, ast.NewSwitchCase("\"list\"", d.desugarListCases(v, cs, dc)))
	}

	if cs, ok := css[dictPattern]; ok {
		ks = append(ks, ast.NewSwitchCase("\"dict\"", d.desugarDictCases(v, cs, dc)))
	}

	if cs, ok := css[scalarPattern]; ok {
		dc = d.desugarScalarCases(v, cs, dc)
	}

	return newSwitch(app("$typeOf", v), ks, dc)
}

func groupCases(cs []ast.MatchCase) map[patternType][]ast.MatchCase {
	css := map[patternType][]ast.MatchCase{}

	for i, c := range cs {
		t := getPatternType(c.Pattern())

		if t == namePattern && i < len(cs)-1 {
			panic("A wildcard pattern is found while some patterns are left")
		}

		css[t] = append(css[t], c)
	}

	return css
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

func (d *desugarer) desugarListCases(v interface{}, cs []ast.MatchCase, dc interface{}) interface{} {
	type group struct {
		first interface{}
		cases []ast.MatchCase
	}

	gs := []group{}

	for i, c := range cs {
		ps := c.Pattern().(ast.App).Arguments().Positionals()

		if len(ps) == 0 {
			dc = app("$if", app("$=", v, "$emptyList"), c.Value(), dc)
			continue
		}

		if ps[0].Expanded() {
			panic("Not implemented")
		}

		first := ps[0].Value()

		c = ast.NewMatchCase(
			ast.NewApp("$list", ast.NewArguments(ps[1:], nil, nil), debug.NewGoInfo(0)),
			c.Value())

		if getPatternType(first) == namePattern {
			d.bindName(first.(string), app("$first", v))
			dc = d.desugarCases(
				app("$rest", v),
				[]ast.MatchCase{c},
				d.desugarListCases(v, cs[i+1:], dc))
			break
		}

		groupExist := false

		for i, g := range gs {
			if equalPatterns(first, g.first) {
				groupExist = true
				gs[i].cases = append(gs[i].cases, c)
			}
		}

		if !groupExist {
			gs = append(gs, group{first, []ast.MatchCase{c}})
		}
	}

	ks := make([]ast.MatchCase, 0, len(gs))

	for _, g := range gs {
		ks = append(ks, ast.NewMatchCase(g.first, d.desugarCases(app("$rest", v), g.cases, dc)))
	}

	return d.desugarCases(app("$first", v), ks, dc)
}

func (d *desugarer) desugarDictCases(v interface{}, cs []ast.MatchCase, dc interface{}) interface{} {
	type group struct {
		key   interface{}
		cases []ast.MatchCase
	}

	gs := []group{}

	for _, c := range cs {
		ps := c.Pattern().(ast.App).Arguments().Positionals()

		if len(ps) == 0 {
			dc = app("$if", app("$=", v, "$emptyDict"), c.Value(), dc)
			continue
		}

		if ps[0].Expanded() {
			panic("Not implemented")
		}

		new := group{ps[0].Value(), []ast.MatchCase{c}}

		if len(gs) == 0 {
			gs = append(gs, new)
		} else if last := gs[len(gs)-1]; equalPatterns(new.key, last.key) {
			last.cases = append(last.cases, c)
		} else {
			gs = append(gs, new)
		}
	}

	x := dc

	for i := len(gs) - 1; i >= 0; i-- {
		g := gs[i]
		x = app("$if",
			app("$include", v, g.key),
			d.desugarDictCasesOfSameKey(v, g.cases, x),
			x)
	}

	return x
}

func (d *desugarer) desugarDictCasesOfSameKey(v interface{}, cs []ast.MatchCase, dc interface{}) interface{} {
	type group struct {
		value interface{}
		cases []ast.MatchCase
	}

	key := cs[0].Pattern().(ast.App).Arguments().Positionals()[0].Value()
	gs := []group{}

	for i, c := range cs {
		ps := c.Pattern().(ast.App).Arguments().Positionals()
		value := ps[1].Value()

		c = ast.NewMatchCase(
			ast.NewApp("$dict", ast.NewArguments(ps[2:], nil, nil), debug.NewGoInfo(0)),
			c.Value())

		if getPatternType(value) == namePattern {
			d.bindName(value.(string), app(v, key))

			if rest := cs[i+1:]; len(rest) != 0 {
				dc = d.desugarDictCasesOfSameKey(v, rest, dc)
			}

			dc = d.desugarCases(app("$delete", v, key), []ast.MatchCase{c}, dc)

			break
		}

		groupExist := false

		for i, g := range gs {
			if equalPatterns(value, g.value) {
				groupExist = true
				gs[i].cases = append(gs[i].cases, c)
			}
		}

		if !groupExist {
			gs = append(gs, group{value, []ast.MatchCase{c}})
		}
	}

	cs = make([]ast.MatchCase, 0, len(gs))

	for _, g := range gs {
		cs = append(
			cs,
			ast.NewMatchCase(g.value, d.desugarCases(app("$delete", v, key), g.cases, dc)))
	}

	return d.desugarCases(app(v, key), cs, dc)
}

func (d *desugarer) desugarScalarCases(v interface{}, cs []ast.MatchCase, dc interface{}) interface{} {
	ks := []ast.SwitchCase{}

	for _, c := range cs {
		ks = append(ks, ast.NewSwitchCase(c.Pattern().(string), c.Value()))
	}

	return newSwitch(v, ks, dc)
}

func renameBoundNamesInCase(c ast.MatchCase) ast.MatchCase {
	p, ns := newPatternRenamer().rename(c.Pattern())
	return ast.NewMatchCase(p, newValueRenamer(ns).rename(c.Value()))
}

func equalPatterns(p, q interface{}) bool {
	switch x := p.(type) {
	case string:
		y, ok := q.(string)

		if !ok {
			return false
		}

		return x == y
	case ast.App:
		y, ok := q.(ast.App)

		if !ok ||
			x.Function().(string) != y.Function().(string) ||
			len(x.Arguments().Positionals()) != len(y.Arguments().Positionals()) {
			return false
		}

		for i := range x.Arguments().Positionals() {
			p := x.Arguments().Positionals()[i]
			q := y.Arguments().Positionals()[i]

			if p.Expanded() != q.Expanded() || !equalPatterns(p.Value(), q.Value()) {
				return false
			}
		}

		return true
	}

	panic(fmt.Errorf("Invalid pattern: %#v, %#v", p, q))
}
