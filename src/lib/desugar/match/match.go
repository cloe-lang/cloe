package match

// Desugar desugars match expressions in a statement.
func Desugar(s interface{}) []interface{} {
	d := newDesugarer()
	l := d.desugarMatchExpression(s)
	return append(d.lets, l)
}
