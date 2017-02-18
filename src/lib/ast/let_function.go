package ast

import "../vm"

type LetFunction struct {
	name      string
	signature vm.Signature
	body      interface{}
}

func NewLetFunction(name string, sig vm.Signature, expr interface{}) LetFunction {
	return LetFunction{name, sig, expr}
}

func (f LetFunction) Name() string {
	return f.name
}

func (f LetFunction) Signature() vm.Signature {
	return f.signature
}

func (f LetFunction) Body() interface{} {
	return f.body
}
