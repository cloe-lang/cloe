package ast

import "fmt"

// LetMatch represents a let-match statement in ASTs.
type LetMatch struct {
	pattern interface{}
	expr    interface{}
}

// NewLetMatch creates a let-match statement from a pattern and an expression.
func NewLetMatch(pattern interface{}, expr interface{}) LetMatch {
	return LetMatch{pattern, expr}
}

// Pattern returns a pattern matched with a value of the let-match statement.
func (m LetMatch) Pattern() interface{} {
	return m.pattern
}

// Expr returns an expression of a variable value""
func (m LetMatch) Expr() interface{} {
	return m.expr
}

func (m LetMatch) String() string {
	return fmt.Sprintf("(let %v %v)\n", m.pattern, m.expr)
}
