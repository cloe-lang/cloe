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
		ps := x.Arguments().Positionals()
		newPs := make([]ast.PositionalArgument, 0, len(ps))

		for _, p := range ps {
			newPs = append(newPs, ast.NewPositionalArgument(r.rename(p.Value()), p.Expanded()))
		}

		ks := x.Arguments().Keywords()
		newKs := make([]ast.KeywordArgument, 0, len(ks))

		for _, k := range ks {
			newKs = append(newKs, ast.NewKeywordArgument(k.Name(), r.rename(k.Value())))
		}

		ds := x.Arguments().ExpandedDicts()
		newDs := make([]interface{}, 0, len(ds))

		for _, d := range ds {
			newDs = append(newDs, r.rename(d))
		}

		return ast.NewApp(
			r.rename(x.Function()),
			ast.NewArguments(newPs, newKs, newDs),
			x.DebugInfo())
	case ast.Switch:
		cs := make([]ast.SwitchCase, 0, len(x.Cases()))

		for _, c := range x.Cases() {
			cs = append(cs, ast.NewSwitchCase(c.Pattern(), r.rename(c.Value())))
		}

		d := interface{}(nil)

		if x.DefaultCase() != nil {
			d = r.rename(x.DefaultCase())
		}

		return ast.NewSwitch(r.rename(x.Value()), cs, d)
	}

	panic(fmt.Errorf("Invalid value: %#v", v))
}
