package ast

type Output struct {
	expr interface{}
}

func NewOutput(expr interface{}) Output {
	return Output{expr}
}

func (o Output) Expr() interface{} {
	return o.expr
}
