package match

import (
	"fmt"

	"github.com/coel-lang/coel/src/lib/ast"
	"github.com/coel-lang/coel/src/lib/consts"
	"github.com/coel-lang/coel/src/lib/debug"
	"github.com/coel-lang/coel/src/lib/gensym"
	"github.com/coel-lang/coel/src/lib/scalar"
)

type desugarer struct {
	letBoundNames, lets []interface{}
}

func newDesugarer() *desugarer {
	return &desugarer{nil, nil}
}

func (d *desugarer) desugar(x interface{}) interface{} {
	return ast.Convert(func(x interface{}) interface{} {
		switch x := x.(type) {
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
		case ast.Match:
			cs := make([]ast.MatchCase, 0, len(x.Cases()))

			for _, c := range x.Cases() {
				cs = append(cs, renameBoundNamesInCase(ast.NewMatchCase(c.Pattern(), d.desugar(c.Value()))))
			}

			return d.resultApp(d.createMatchFunction(cs), d.desugar(x.Value()))
		}

		return nil
	}, x)
}

func (d *desugarer) takeLets() []interface{} {
	ls := append(d.letBoundNames, d.lets...)
	d.letBoundNames = nil
	d.lets = nil
	return ls
}

func (d *desugarer) letTempVar(v interface{}) string {
	s := gensym.GenSym()
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
	return d.bindName(gensym.GenSym(), app(f, args...))
}

// resultApp applies a function to arguments and creates a result value of match
// expression.
func (d *desugarer) resultApp(f interface{}, args ...interface{}) string {
	return d.letTempVar(app(f, args...))
}

func (d *desugarer) createMatchFunction(cs []ast.MatchCase) interface{} {
	arg := gensym.GenSym()
	body := d.desugarCases(arg, cs, "$matchError")

	f := ast.NewLetFunction(
		gensym.GenSym(),
		ast.NewSignature([]string{arg}, nil, "", nil, nil, ""),
		d.takeLets(),
		app("$if", app("$=", app("$catch", arg), "$nil"), body, arg),
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
		case consts.Names.ListFunction:
			return listPattern
		case consts.Names.DictionaryFunction:
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
		case consts.Names.ListFunction:
			return ok
		case consts.Names.DictionaryFunction:
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
	emptyCase := (interface{})(nil)

	for i, c := range cs {
		ps := c.Pattern().(ast.App).Arguments().Positionals()

		if len(ps) == 0 {
			emptyCase = c.Value()
			continue
		}

		v := ps[0].Value()

		if ps[0].Expanded() {
			d.bindName(v.(string), list)
			dc = c.Value()
			break
		}

		c = ast.NewMatchCase(
			ast.NewApp(
				consts.Names.ListFunction,
				ast.NewArguments(ps[1:], nil, nil),
				debug.NewGoInfo(0)),
			c.Value())

		if isGeneralNamePattern(v) {
			d.bindName(v, first)

			if cs := cs[i+1:]; len(cs) > 0 {
				dc = d.desugarListCases(list, cs, dc)
			}

			dc = d.defaultCaseOfGeneralNamePattern(
				first,
				v,
				d.desugarListCases(rest, []ast.MatchCase{c}, dc),
				dc)
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

	r := d.desugarCases(first, cs, dc)

	if emptyCase == nil {
		return r
	}

	return d.resultApp("$if", app("$=", list, consts.Names.EmptyList), emptyCase, r)
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
			dc = d.resultApp("$if", app("$=", dict, consts.Names.EmptyDictionary), c.Value(), dc)
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
	rest := d.matchedApp("delete", dict, key)
	gs := []group{}

	for i, c := range cs {
		ps := c.Pattern().(ast.App).Arguments().Positionals()
		v := ps[1].Value()

		c = ast.NewMatchCase(
			ast.NewApp(
				consts.Names.DictionaryFunction,
				ast.NewArguments(ps[2:], nil, nil),
				debug.NewGoInfo(0)),
			c.Value())

		if isGeneralNamePattern(v) {
			d.bindName(v, value)

			if cs := cs[i+1:]; len(cs) != 0 {
				dc = d.desugarDictCasesOfSameKey(dict, cs, dc)
			}

			dc = d.defaultCaseOfGeneralNamePattern(
				value,
				v,
				d.desugarCases(rest, []ast.MatchCase{c}, dc),
				dc)
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
		cs = append(cs, ast.NewMatchCase(g.value, d.desugarCases(rest, g.cases, dc)))
	}

	return d.desugarCases(value, cs, dc)
}

func (d *desugarer) defaultCaseOfGeneralNamePattern(v, p, body, dc interface{}) interface{} {
	switch getPatternType(p) {
	case namePattern:
		return body
	case listPattern:
		return d.ifType(v, "list", body, dc)
	case dictPattern:
		return d.ifType(v, "dict", body, dc)
	}

	panic("Unreachable")
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
