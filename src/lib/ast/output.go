package ast

type Output struct {
	expr     interface{}
	expanded bool
}

func NewOutput(expr interface{}, expanded bool) Output {
	return Output{expr, expanded}
}

func (o Output) Expr() interface{} {
	return o.expr
}
