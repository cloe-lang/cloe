package desugar

import "github.com/raviqqe/tisp/src/lib/ast"

type desugarer struct {
	statements []interface{}
}

func newDesugarer() desugarer {
	return desugarer{}
}

func (d *desugarer) desugar(module []interface{}) []interface{} {
	d.statements = make([]interface{}, 0, 2*len(module)) // TODO: Best cap?

	for _, s := range module {
		d.appendStatement(d.desugarStatement(s))
	}

	return d.statements
}

func (d *desugarer) appendStatement(s interface{}) {
	d.statements = append(d.statements, s)
}

func (d *desugarer) desugarStatement(s interface{}) interface{} {
	switch s := s.(type) {
	case ast.LetFunction:
		return s
	default:
		return s
	}
}
