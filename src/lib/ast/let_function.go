package ast

type LetFunction struct {
	name      string
	signature Signature
	lets      []interface{}
	body      interface{}
}

func NewLetFunction(name string, sig Signature, lets []interface{}, expr interface{}) LetFunction {
	return LetFunction{name, sig, lets, expr}
}

func (f LetFunction) Name() string {
	return f.name
}

func (f LetFunction) Signature() Signature {
	return f.signature
}

// Lets returns let statements contained in this let-function statement.
// Return values should be LetConst or LetFunction.
func (f LetFunction) Lets() []interface{} {
	return f.lets
}

func (f LetFunction) Body() interface{} {
	return f.body
}
