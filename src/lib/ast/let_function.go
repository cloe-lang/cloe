package ast

type LetFunction struct {
	name      string
	signature Signature
	body      interface{}
}

func NewLetFunction(name string, sig Signature, expr interface{}) LetFunction {
	return LetFunction{name, sig, expr}
}

func (f LetFunction) Name() string {
	return f.name
}

func (f LetFunction) Signature() Signature {
	return f.signature
}

func (f LetFunction) Body() interface{} {
	return f.body
}
