package match

import (
	"fmt"

	"github.com/cloe-lang/cloe/src/lib/ast"
	"github.com/cloe-lang/cloe/src/lib/consts"
	"github.com/cloe-lang/cloe/src/lib/debug"
	"github.com/cloe-lang/cloe/src/lib/gensym"
	"github.com/cloe-lang/cloe/src/lib/scalar"
)

type casesDesugarer struct {
	letBoundNames, lets []interface{}
}

func newCasesDesugarer() *casesDesugarer {
	return &casesDesugarer{nil, nil}
}

func (d *casesDesugarer) Desugar(cs []ast.MatchCase) ast.DefFunction {
	arg := gensym.GenSym()
	body := d.desugarCases(arg, cs, "$matchError")

	return ast.NewDefFunction(
		gensym.GenSym(),
		ast.NewSignature([]string{arg}, "", nil, ""),
		d.takeLets(),
		app("$if", app("$=", app("$catch", arg), "$nil"), body, arg),
		debug.NewGoInfo(0))
}

func (d *casesDesugarer) takeLets() []interface{} {
	ls := append(d.letBoundNames, d.lets...)
	d.letBoundNames = nil
	d.lets = nil
	return ls
}

func (d *casesDesugarer) letTempVar(v interface{}) string {
	s := gensym.GenSym()
	d.lets = append(d.lets, ast.NewLetVar(s, v))
	return s
}

func (d *casesDesugarer) bindName(p interface{}, v interface{}) string {
	s := generalNamePatternToName(p)
	d.letBoundNames = append(d.letBoundNames, ast.NewLetVar(s, v))
	return s
}

// matchedApp applies a function to arguments and creates a matched value of
// match expression.
func (d *casesDesugarer) matchedApp(f interface{}, args ...interface{}) string {
	return d.bindName(gensym.GenSym(), app(f, args...))
}

// resultApp applies a function to arguments and creates a result value of match
// expression.
func (d *casesDesugarer) resultApp(f interface{}, args ...interface{}) string {
	return d.letTempVar(app(f, args...))
}

func (d *casesDesugarer) desugarCases(v interface{}, cs []ast.MatchCase, dc interface{}) interface{} {
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

	if cs, ok := css[dictionaryPattern]; ok {
		ks = append(ks, ast.NewSwitchCase("\"dictionary\"", d.desugarDictionaryCases(v, cs, dc)))
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
			return dictionaryPattern
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

func (d *casesDesugarer) desugarListCases(list interface{}, cs []ast.MatchCase, dc interface{}) interface{} {
	type group struct {
		first interface{}
		cases []ast.MatchCase
	}

	gs := []group{}
	first := d.matchedApp("$first", list)
	rest := d.matchedApp("$rest", list)
	emptyCase := interface{}(nil)

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
				ast.NewArguments(ps[1:], nil),
				debug.NewGoInfo(0)),
			c.Value())

		var ok bool
		if dc, ok = d.handleGeneralNamePattern(
			v, first, cs, c, i, dc, list, rest, d.desugarListCases, d.desugarListCases); ok {
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

func (d *casesDesugarer) ifType(v interface{}, t string, then, els interface{}) interface{} {
	return d.resultApp("$if", app("$=", app("$typeOf", v), "\""+t+"\""), then, els)
}

func (d *casesDesugarer) desugarDictionaryCases(dict interface{}, cs []ast.MatchCase, dc interface{}) interface{} {
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
			d.desugarDictionaryCasesOfSameKey(dict, g.cases, dc),
			dc)
	}

	return dc
}

func (d *casesDesugarer) desugarDictionaryCasesOfSameKey(dict interface{}, cs []ast.MatchCase, dc interface{}) interface{} {
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
				ast.NewArguments(ps[2:], nil),
				debug.NewGoInfo(0)),
			c.Value())

		var ok bool
		if dc, ok = d.handleGeneralNamePattern(
			v, value, cs, c, i, dc, dict, rest, d.desugarDictionaryCasesOfSameKey, d.desugarDictionaryCases); ok {
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

func (d *casesDesugarer) handleGeneralNamePattern(
	p, v interface{}, cs []ast.MatchCase, c ast.MatchCase, i int,
	dc, original, rest interface{},
	desugarRestCases, desugarCollectionCases func(interface{}, []ast.MatchCase, interface{}) interface{},
) (interface{}, bool) {
	if isGeneralNamePattern(p) {
		d.bindName(p, v)

		if cs := cs[i+1:]; len(cs) > 0 {
			dc = desugarRestCases(original, cs, dc)
		}

		dc = d.defaultCaseOfGeneralNamePattern(
			v,
			p,
			desugarCollectionCases(rest, []ast.MatchCase{c}, dc),
			dc)

		return dc, true
	}

	return dc, false
}

func (d *casesDesugarer) defaultCaseOfGeneralNamePattern(v, p, body, dc interface{}) interface{} {
	switch getPatternType(p) {
	case namePattern:
		return body
	case listPattern:
		return d.ifType(v, "list", body, dc)
	case dictionaryPattern:
		return d.ifType(v, "dictionary", body, dc)
	}

	panic("Unreachable")
}

func (d *casesDesugarer) desugarScalarCases(v interface{}, cs []ast.MatchCase, dc interface{}) interface{} {
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
