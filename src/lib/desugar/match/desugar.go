package match

import "github.com/cloe-lang/cloe/src/lib/ast"

// Desugar desugars match expressions in a statement.
func Desugar(s interface{}) []interface{} {
	return []interface{}{desugarMatch(s)}
}

func desugarMatch(x interface{}) interface{} {
	return ast.Convert(func(x interface{}) interface{} {
		m, ok := x.(ast.Match)

		if !ok {
			return nil
		}

		cs := make([]ast.MatchCase, 0, len(m.Cases()))

		for _, c := range m.Cases() {
			p, r := newPatternRenamer().Rename(c.Pattern())
			cs = append(cs, ast.NewMatchCase(p, r.Rename(desugarMatch(c.Value()))))
		}

		return app(newCasesDesugarer().Desugar(cs), desugarMatch(m.Value()))
	}, x)
}
