package ast

type LetVar struct {
	name string
	expr interface{}
}

func NewLetVar(name string, expr interface{}) LetVar {
	return LetVar{name, expr}
}

func (c LetVar) Name() string {
	return c.name
}

func (c LetVar) Expr() interface{} {
	return c.expr
}
