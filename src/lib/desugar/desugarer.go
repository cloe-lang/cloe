package desugar

type desugarer struct{}

func newDesugarer() desugarer {
	return desugarer{}
}

func (d *desugarer) desugar(module []interface{}) []interface{} {
	ss := make([]interface{}, 0, 2*len(module)) // TODO: Best cap?

	for _, s := range module {
		ss = append(ss, d.desugarStatement(s))
	}

	return ss
}

func (*desugarer) desugarStatement(s interface{}) interface{} {
	switch s := s.(type) {
	default:
		return s
	}
}
