package desugar

func Desugar(module []interface{}) []interface{} {
	d := newDesugarer()
	return d.desugar(module)
}
