package ast

import "fmt"

// AnonymousFunction represents a anonymous function as an expression.
type AnonymousFunction struct {
	signature Signature
	lets      []any
	body      any
}

// NewAnonymousFunction creates a anonymous function.
func NewAnonymousFunction(s Signature, ls []any, b any) AnonymousFunction {
	return AnonymousFunction{s, ls, b}
}

// Signature returns a signature of an anonymous function.
func (f AnonymousFunction) Signature() Signature {
	return f.signature
}

// Lets returns let statements contained in an anonymous function.
func (f AnonymousFunction) Lets() []any {
	return f.lets
}

// Body returns a body expression of an anonymous function.
func (f AnonymousFunction) Body() any {
	return f.body
}

func (f AnonymousFunction) String() string {
	return fmt.Sprintf("(\\ (%v) %v)", f.signature, f.body)
}
