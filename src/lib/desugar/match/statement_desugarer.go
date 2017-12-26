package match

import (
	"github.com/coel-lang/coel/src/lib/ast"
	"github.com/coel-lang/coel/src/lib/gensym"
)

type statementDesugarer struct {
	lets []interface{}
}

func newStatementDesugarer() *statementDesugarer {
	return &statementDesugarer{nil}
}

func (d *statementDesugarer) Desugar(x interface{}) interface{} {
	return ast.Convert(func(x interface{}) interface{} {
		switch x := x.(type) {
		case ast.DefFunction:
			ls := make([]interface{}, 0, len(x.Lets()))

			for _, l := range x.Lets() {
				l := d.Desugar(l)
				ls = append(ls, append(d.TakeLets(), l)...)
			}

			b := d.Desugar(x.Body())

			return ast.NewDefFunction(
				x.Name(),
				x.Signature(),
				append(ls, d.TakeLets()...),
				b,
				x.DebugInfo())
		case ast.Match:
			cs := make([]ast.MatchCase, 0, len(x.Cases()))

			for _, c := range x.Cases() {
				cs = append(cs, renameBoundNamesInCase(ast.NewMatchCase(c.Pattern(), c.Value())))
			}

			f := newCasesDesugarer().Desugar(cs)
			s := gensym.GenSym()

			d.lets = append(d.lets, append(
				Desugar(f),
				ast.NewLetVar(s, app(f.Name(), d.Desugar(x.Value()))))...)

			return s
		}

		return nil
	}, x)
}

func (d *statementDesugarer) TakeLets() []interface{} {
	ls := d.lets
	d.lets = nil
	return ls
}
