package ast

import "fmt"

// Effect represents effects of programs.
type Effect struct {
	expr     any
	expanded bool
}

// NewEffect creates an Effect.
func NewEffect(expr any, expanded bool) Effect {
	return Effect{expr, expanded}
}

// Expr returns an expression of the effect.
func (o Effect) Expr() any {
	return o.expr
}

// Expanded returns true when the effect is expanded.
func (o Effect) Expanded() bool {
	return o.expanded
}

func (o Effect) String() string {
	if o.expanded {
		return fmt.Sprintf("..%v", o.expr)
	}

	return fmt.Sprintf("%v", o.expr)
}
