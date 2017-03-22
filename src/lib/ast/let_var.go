package ast

// LetVar represents a let-variable statement node in ASTs.
type LetVar struct {
	name string
	expr interface{}
}

// NewLetVar creates a LetVar from a variable name and its value of an expression.
func NewLetVar(name string, expr interface{}) LetVar {
	return LetVar{name, expr}
}

// Name returns a variable name defined by the let-variable statement.
func (c LetVar) Name() string {
	return c.name
}

// Expr returns an expression of a variable value""
func (c LetVar) Expr() interface{} {
	return c.expr
}
