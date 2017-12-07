package ast

import (
	"fmt"
	"strings"

	"github.com/kr/text"
	"github.com/tisp-lang/tisp/src/lib/debug"
)

// LetFunction represents a let-function statement node in ASTs.
type LetFunction struct {
	name      string
	signature Signature
	lets      []interface{}
	body      interface{}
	info      debug.Info
}

// NewLetFunction creates a LetFunction from its function name, signature,
// internal let statements, and a body expression.
func NewLetFunction(name string, sig Signature, lets []interface{}, expr interface{}, i debug.Info) LetFunction {
	return LetFunction{name, sig, lets, expr, i}
}

// Name returns a name of a function defined by the let-function statement.
func (f LetFunction) Name() string {
	return f.name
}

// Signature returns a signature of a function defined by the let-function statement.
func (f LetFunction) Signature() Signature {
	return f.signature
}

// Lets returns let statements contained in the let-function statement.
// Returned values should be LetVar or LetFunction.
func (f LetFunction) Lets() []interface{} {
	return f.lets
}

// Body returns a body expression of a function defined by the let-function statement.
func (f LetFunction) Body() interface{} {
	return f.body
}

// DebugInfo returns debug information of a function defined by the let-function statement.
func (f LetFunction) DebugInfo() debug.Info {
	return f.info
}

func (f LetFunction) String() string {
	ss := make([]string, 0, len(f.lets))

	for _, l := range f.lets {
		ss = append(ss, text.Indent(l.(fmt.Stringer).String(), "\t"))
	}

	return fmt.Sprintf("(def (%v %v)\n%v\t%v)\n", f.name, f.signature, strings.Join(ss, ""), f.body)
}
