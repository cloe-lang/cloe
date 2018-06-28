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
		m := map[string]string{}
		ns := x.Signature().NameToIndex()

		for k, v := range r.nameMap {
			if _, ok := ns[k]; !ok {
				m[k] = v
			}
		}

		return ast.NewAnonymousFunction(x.Signature(), newValueRenamer(m).rename(x.Body()))
	case ast.Match:
		cs := make([]ast.MatchCase, 0, len(x.Cases()))

		for _, c := range x.Cases() {
			p, ns := newPatternRenamer().rename(c.Pattern())
			cs = append(cs, ast.NewMatchCase(p, r.extend(ns).rename(c.Value())))
		}

		return ast.NewMatch(r.rename(x.Value()), cs)
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
