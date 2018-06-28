package match

import (
	"fmt"

	"github.com/cloe-lang/cloe/src/lib/ast"
)

type valueRenamer struct {
	nameMap map[string]string
}

func newValueRenamer(m map[string]string) valueRenamer {
	return valueRenamer{m}
}

func (r valueRenamer) rename(v interface{}) interface{} {
	switch x := v.(type) {
	case string:
		if n, ok := r.nameMap[x]; ok {
			return n
		}

		return x
	case ast.App:
		args := x.Arguments()
		ps := make([]ast.PositionalArgument, 0, len(args.Positionals()))

		for _, p := range args.Positionals() {
			ps = append(ps, ast.NewPositionalArgument(r.rename(p.Value()), p.Expanded()))
		}

		ks := make([]ast.KeywordArgument, 0, len(args.Keywords()))

		for _, k := range args.Keywords() {
			ks = append(ks, ast.NewKeywordArgument(k.Name(), r.rename(k.Value())))
		}

		return ast.NewApp(r.rename(x.Function()), ast.NewArguments(ps, ks), x.DebugInfo())
	case ast.AnonymousFunction:
		r := r.copy()

		for n := range x.Signature().NameToIndex() {
			r.delete(n)
		}

		ls := make([]interface{}, 0, len(x.Lets()))

		for _, l := range x.Lets() {
			switch x := l.(type) {
			case ast.LetVar:
				ls = append(ls, ast.NewLetVar(x.Name(), r.rename(x.Expr())))
				r.delete(x.Name())
			default:
				panic("unreachable")
			}
		}

		return ast.NewAnonymousFunction(x.Signature(), ls, r.rename(x.Body()))
	case ast.Switch:
		cs := make([]ast.SwitchCase, 0, len(x.Cases()))

		for _, c := range x.Cases() {
			// Switch case patterns are already renamed as match case patterns.
			cs = append(cs, ast.NewSwitchCase(c.Pattern(), r.rename(c.Value())))
		}

		return ast.NewSwitch(r.rename(x.Value()), cs, r.rename(x.DefaultCase()))
	}

	panic(fmt.Errorf("Invalid value: %#v", v))
}

func (r valueRenamer) extend(ns map[string]string) valueRenamer {
	ms := make(map[string]string, len(r.nameMap)+len(ns))

	for _, ns := range []map[string]string{r.nameMap, ns} {
		for n, m := range ns {
			ms[n] = m
		}
	}

	return newValueRenamer(ms)
}

func (r valueRenamer) copy() valueRenamer {
	m := make(map[string]string, len(r.nameMap))

	for k, v := range r.nameMap {
		m[k] = v
	}

	return newValueRenamer(m)
}

func (r valueRenamer) delete(s string) {
	delete(r.nameMap, s)
}
