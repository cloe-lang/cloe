package ast

// AnonymousFunction represents a anonymous function as an expression.
type AnonymousFunction struct {
	signature Signature
	body      interface{}
}

// NewAnonymousFunction creates a anonymous function.
func NewAnonymousFunction(s Signature, b interface{}) AnonymousFunction {
	return AnonymousFunction{s, b}
}

// Signature returns a signature of an anonymous function.
func (f AnonymousFunction) Signature() Signature {
	return f.signature
}

// Body returns a body expression of an anonymous function.
func (f AnonymousFunction) Body() interface{} {
	return f.body
}
