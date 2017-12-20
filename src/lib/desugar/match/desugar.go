package match

// Desugar desugars match expressions in a statement.
func Desugar(s interface{}) []interface{} {
	d := newStatementDesugarer()
	l := d.Desugar(s)
	return append(d.TakeLets(), l)
}
