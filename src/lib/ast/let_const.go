package ast

type LetConst struct {
	name string
	expr interface{}
}

func NewLetConst(name string, expr interface{}) LetConst {
	return LetConst{name, expr}
}

func (c LetConst) Name() string {
	return c.name
}

func (c LetConst) Expr() interface{} {
	return c.expr
}
