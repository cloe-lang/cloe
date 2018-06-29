package ast

import (
	"fmt"
	"strings"
)

// LetExpression represents a let expression node in ASTs.
type LetExpression struct {
	lets []interface{}
	expr interface{}
}

// NewLetExpression creates a let expression.
func NewLetExpression(lets []interface{}, expr interface{}) LetExpression {
	return LetExpression{lets, expr}
}

// Lets returns let-match statements used in the let expression.
func (l LetExpression) Lets() []interface{} {
	return l.lets
}

// Expr returns an expression of the let expression.
func (l LetExpression) Expr() interface{} {
	return l.expr
}

func (l LetExpression) String() string {
	ss := []string{}

	for _, l := range l.lets {
		switch l := l.(type) {
		case LetVar:
			ss = append(ss, l.Name(), fmt.Sprint(l.Expr()))
		case LetMatch:
			ss = append(ss, fmt.Sprint(l.Pattern(), l.Expr()))
		default:
			panic("unreachable")
		}
	}

	return fmt.Sprintf("(let %v %v)", strings.Join(ss, " "), l.expr)
}
