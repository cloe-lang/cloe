package desugar

import "github.com/tisp-lang/tisp/src/lib/ast"

func desugarDictionaryExpansion(x interface{}) []interface{} {
	return []interface{}{ast.Convert(func(x interface{}) interface{} {
		a, ok := x.(ast.App)

		if !ok || a.Function() != "$dict" {
			return nil
		}

		ps := a.Arguments().Positionals()
		args := make([]interface{}, 0, len(ps))
		dicts := make([]interface{}, 0, len(ps))

		for _, p := range ps {
			if p.Expanded() {
				dicts = append(dicts, p.Value())
			} else {
				args = append(args, p.Value())
			}
		}

		return ast.NewPApp("$merge", append([]interface{}{ast.NewPApp(a.Function(), args, a.DebugInfo())}, dicts...), a.DebugInfo())
	}, x)}
}
