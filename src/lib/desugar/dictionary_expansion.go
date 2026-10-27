package desugar

import (
	"github.com/cloe-lang/cloe/src/lib/ast"
	"github.com/cloe-lang/cloe/src/lib/consts"
)

func desugarDictionaryExpansion(x any) []any {
	return []any{ast.Convert(func(x any) any {
		a, ok := x.(ast.App)

		if !ok || a.Function() != consts.Names.DictionaryFunction {
			return nil
		}

		ps := a.Arguments().Positionals()
		args := make([]any, 0, len(ps))
		dicts := make([]any, 0, len(ps))

		for _, p := range ps {
			if p.Expanded() {
				dicts = append(dicts, p.Value())
			} else {
				args = append(args, p.Value())
			}
		}

		return ast.NewPApp("$merge", append([]any{ast.NewPApp(a.Function(), args, a.DebugInfo())}, dicts...), a.DebugInfo())
	}, x)}
}
