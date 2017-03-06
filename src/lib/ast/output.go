package ast

// Output represents outputs of programs.
type Output struct {
	expr     interface{}
	expanded bool
}

// NewOutput creates an Output.
func NewOutput(expr interface{}, expanded bool) Output {
	return Output{expr, expanded}
}

// Expr returns an expression of the output.
func (o Output) Expr() interface{} {
	return o.expr
}

// Expanded returns true when the output is expanded.
func (o Output) Expanded() bool {
	return o.expanded
}
