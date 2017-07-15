package ast

import "fmt"

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
func (v LetVar) Name() string {
	return v.name
}

// Expr returns an expression of a variable value""
func (v LetVar) Expr() interface{} {
	return v.expr
}

func (v LetVar) String() string {
	return fmt.Sprintf("(let %v %v)", v.name, v.expr)
}
