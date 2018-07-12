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

func (r valueRenamer) Rename(v interface{}) interface{} {
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
			ps = append(ps, ast.NewPositionalArgument(r.Rename(p.Value()), p.Expanded()))
		}

		ks := make([]ast.KeywordArgument, 0, len(args.Keywords()))

		for _, k := range args.Keywords() {
			ks = append(ks, ast.NewKeywordArgument(k.Name(), r.Rename(k.Value())))
		}

		return ast.NewApp(r.Rename(x.Function()), ast.NewArguments(ps, ks), x.DebugInfo())
	case ast.AnonymousFunction:
		r := r.copy()

		for _, n := range x.Signature().Names() {
			r.delete(n)
		}

		ls := make([]interface{}, 0, len(x.Lets()))

		for _, l := range x.Lets() {
			switch x := l.(type) {
			case ast.LetVar:
				ls = append(ls, ast.NewLetVar(x.Name(), r.Rename(x.Expr())))
				r.delete(x.Name())
			default:
				panic("unreachable")
			}
		}

		return ast.NewAnonymousFunction(x.Signature(), ls, r.Rename(x.Body()))
	case ast.Switch:
		cs := make([]ast.SwitchCase, 0, len(x.Cases()))

		for _, c := range x.Cases() {
			// Switch case patterns are already renamed as match case patterns.
			cs = append(cs, ast.NewSwitchCase(c.Pattern(), r.Rename(c.Value())))
		}

		return ast.NewSwitch(r.Rename(x.Value()), cs, r.Rename(x.DefaultCase()))
	}

	panic(fmt.Errorf("Invalid value: %#v", v))
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
