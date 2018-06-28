package ast

import "fmt"

// AnonymousFunction represents a anonymous function as an expression.
type AnonymousFunction struct {
	signature Signature
	lets      []interface{}
	body      interface{}
}

// NewAnonymousFunction creates a anonymous function.
func NewAnonymousFunction(s Signature, ls []interface{}, b interface{}) AnonymousFunction {
	return AnonymousFunction{s, ls, b}
}

// Signature returns a signature of an anonymous function.
func (f AnonymousFunction) Signature() Signature {
	return f.signature
}

// Lets returns let statements contained in an anonymous function.
func (f AnonymousFunction) Lets() []interface{} {
	return f.lets
}

// Body returns a body expression of an anonymous function.
func (f AnonymousFunction) Body() interface{} {
	return f.body
}

func (f AnonymousFunction) String() string {
	return fmt.Sprintf("(\\ (%v) %v)", f.signature, f.body)
}
