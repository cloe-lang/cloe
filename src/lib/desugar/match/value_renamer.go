package match

import (
	"fmt"

	"github.com/tisp-lang/tisp/src/lib/ast"
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

		ds := make([]interface{}, 0, len(args.ExpandedDicts()))

		for _, d := range args.ExpandedDicts() {
			ds = append(ds, r.rename(d))
		}

		return ast.NewApp(r.rename(x.Function()), ast.NewArguments(ps, ks, ds), x.DebugInfo())
	case ast.Switch:
		cs := make([]ast.SwitchCase, 0, len(x.Cases()))

		for _, c := range x.Cases() {
			cs = append(cs, ast.NewSwitchCase(c.Pattern(), r.rename(c.Value())))
		}

		return newSwitch(r.rename(x.Value()), cs, r.rename(x.DefaultCase()))
	}

	panic(fmt.Errorf("Invalid value: %#v", v))
}
