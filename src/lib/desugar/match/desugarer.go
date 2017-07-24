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

		ds := make([]interface{}, 0, len(x.ExpandedDicts()))

		for _, dict := range x.ExpandedDicts() {
			ds = append(ds, d.desugar(dict))
		}

		return ast.NewArguments(ps, ks, ds)
	case ast.KeywordArgument:
		return ast.NewKeywordArgument(x.Name(), d.desugar(x.Value()))
	case ast.LetFunction:
		ls := make([]interface{}, 0, len(x.Lets()))

		for _, l := range x.Lets() {
			l := d.desugar(l)
			ls = append(ls, append(d.takeLets(), l)...)
		}

		b := d.desugar(x.Body())

		return ast.NewLetFunction(
			x.Name(),
			x.Signature(),
			append(ls, d.takeLets()...),
			b,
			x.DebugInfo())
	case ast.LetVar:
		return ast.NewLetVar(x.Name(), d.desugar(x.Expr()))
	case ast.Match:
		cs := make([]ast.MatchCase, 0, len(x.Cases()))

		for _, c := range x.Cases() {
			cs = append(cs, renameBoundNamesInCase(ast.NewMatchCase(c.Pattern(), d.desugar(c.Value()))))
		}

		return d.resultApp(d.createMatchFunction(cs), d.desugar(x.Value()))
	case ast.MutualRecursion:
		fs := make([]ast.LetFunction, 0, len(x.LetFunctions()))

		for _, f := range x.LetFunctions() {
			fs = append(fs, d.desugar(f).(ast.LetFunction))
		}

		return ast.NewMutualRecursion(fs, x.DebugInfo())
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

func (d *desugarer) letTempVar(v interface{}) string {
	s := gensym.GenSym("match", "tmp")
	d.lets = append(d.lets, ast.NewLetVar(s, v))
	return s
}

func (d *desugarer) bindName(p interface{}, v interface{}) string {
	s := generalNamePatternToName(p)
	d.letBoundNames = append(d.letBoundNames, ast.NewLetVar(s, v))
	return s
}

// matchedApp applies a function to arguments and creates a matched value of
// match expression.
func (d *desugarer) matchedApp(f interface{}, args ...interface{}) string {
	return d.bindName(gensym.GenSym("match", "app"), app(f, args...))
}

// resultApp applies a function to arguments and creates a result value of match
// expression.
func (d *desugarer) resultApp(f interface{}, args ...interface{}) string {
	return d.letTempVar(app(f, args...))
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

	return newSwitch(d.resultApp("$typeOf", v), ks, dc)
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

func isGeneralNamePattern(p interface{}) bool {
	switch x := p.(type) {
	case string:
		if scalar.Defined(x) {
			return false
		}

		return true
	case ast.App:
		ps := x.Arguments().Positionals()
		ok := len(ps) == 1 && ps[0].Expanded()

		switch x.Function().(string) {
		case "$list":
			return ok
		case "$dict":
			return ok
		}
	}

	panic(fmt.Errorf("Invalid pattern: %#v", p))
}

func (d *desugarer) desugarListCases(list interface{}, cs []ast.MatchCase, dc interface{}) interface{} {
	type group struct {
		first interface{}
		cases []ast.MatchCase
	}

	gs := []group{}
	first := d.matchedApp("$first", list)
	rest := d.matchedApp("$rest", list)

	for i, c := range cs {
		ps := c.Pattern().(ast.App).Arguments().Positionals()

		if len(ps) == 0 {
			dc = d.resultApp("$if", app("$=", list, "$emptyList"), c.Value(), dc)
			continue
		}

		v := ps[0].Value()

		if ps[0].Expanded() {
			d.bindName(v.(string), list)
			dc = c.Value()
			break
		}

		c = ast.NewMatchCase(
			ast.NewApp("$list", ast.NewArguments(ps[1:], nil, nil), debug.NewGoInfo(0)),
			c.Value())

		if isGeneralNamePattern(v) {
			d.bindName(v, first)

			if cs := cs[i+1:]; len(cs) > 0 {
				dc = d.desugarListCases(list, cs, dc)
			}

			next := d.desugarCases(rest, []ast.MatchCase{c}, dc)

			switch getPatternType(v) {
			case namePattern:
				dc = next
			case listPattern:
				dc = d.ifType(first, "list", next, dc)
			case dictPattern:
				dc = d.ifType(first, "dict", next, dc)
			default:
				panic("Unreachable")
			}

			break
		}

		groupExist := false

		for i, g := range gs {
			if equalPatterns(v, g.first) {
				groupExist = true
				gs[i].cases = append(gs[i].cases, c)
			}
		}

		if !groupExist {
			gs = append(gs, group{v, []ast.MatchCase{c}})
		}
	}

	cs = make([]ast.MatchCase, 0, len(gs))
	dc = d.letTempVar(dc)

	for _, g := range gs {
		cs = append(cs, ast.NewMatchCase(g.first, d.desugarCases(rest, g.cases, dc)))
	}

	return d.desugarCases(first, cs, dc)
}

func (d *desugarer) ifType(v interface{}, t string, then, els interface{}) interface{} {
	return d.resultApp("$if", app("$=", app("$typeOf", v), "\""+t+"\""), then, els)
}

func (d *desugarer) desugarDictCases(dict interface{}, cs []ast.MatchCase, dc interface{}) interface{} {
	type group struct {
		key   interface{}
		cases []ast.MatchCase
	}

	gs := []group{}

	for _, c := range cs {
		ps := c.Pattern().(ast.App).Arguments().Positionals()

		if len(ps) == 0 {
			dc = d.resultApp("$if", app("$=", dict, "$emptyDict"), c.Value(), dc)
			continue
		}

		k := ps[0].Value()

		if ps[0].Expanded() {
			d.bindName(k.(string), dict)
			dc = c.Value()
			break
		}

		g := group{k, []ast.MatchCase{c}}

		if len(gs) == 0 {
			gs = append(gs, g)
		} else if last := len(gs) - 1; equalPatterns(g.key, gs[last].key) {
			gs[last].cases = append(gs[last].cases, c)
		} else {
			gs = append(gs, g)
		}
	}

	for i := len(gs) - 1; i >= 0; i-- {
		g := gs[i]
		dc = d.resultApp("$if",
			app("$include", dict, g.key),
			d.desugarDictCasesOfSameKey(dict, g.cases, dc),
			dc)
	}

	return dc
}

func (d *desugarer) desugarDictCasesOfSameKey(dict interface{}, cs []ast.MatchCase, dc interface{}) interface{} {
	type group struct {
		value interface{}
		cases []ast.MatchCase
	}

	key := cs[0].Pattern().(ast.App).Arguments().Positionals()[0].Value()
	value := d.matchedApp(dict, key)
	newDict := d.matchedApp("delete", dict, key)
	gs := []group{}

	for i, c := range cs {
		ps := c.Pattern().(ast.App).Arguments().Positionals()
		v := ps[1].Value()

		c = ast.NewMatchCase(
			ast.NewApp("$dict", ast.NewArguments(ps[2:], nil, nil), debug.NewGoInfo(0)),
			c.Value())

		if isGeneralNamePattern(v) {
			d.bindName(v, value)

			if cs := cs[i+1:]; len(cs) != 0 {
				dc = d.desugarDictCasesOfSameKey(dict, cs, dc)
			}

			next := d.desugarCases(newDict, []ast.MatchCase{c}, dc)

			switch getPatternType(v) {
			case namePattern:
				dc = next
			case listPattern:
				dc = d.ifType(value, "list", next, dc)
			case dictPattern:
				dc = d.ifType(value, "dict", next, dc)
			default:
				panic("Unreachable")
			}

			break
		}

		groupExist := false

		for i, g := range gs {
			if equalPatterns(v, g.value) {
				groupExist = true
				gs[i].cases = append(gs[i].cases, c)
			}
		}

		if !groupExist {
			gs = append(gs, group{v, []ast.MatchCase{c}})
		}
	}

	cs = make([]ast.MatchCase, 0, len(gs))
	dc = d.letTempVar(dc)

	for _, g := range gs {
		cs = append(cs, ast.NewMatchCase(g.value, d.desugarCases(newDict, g.cases, dc)))
	}

	return d.desugarCases(value, cs, dc)
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

func generalNamePatternToName(p interface{}) string {
	switch x := p.(type) {
	case string:
		return x
	case ast.App:
		if ps := x.Arguments().Positionals(); len(ps) == 1 && ps[0].Expanded() {
			return ps[0].Value().(string)
		}
	}

	panic(fmt.Errorf("Invalid pattern: %#v", p))
}
