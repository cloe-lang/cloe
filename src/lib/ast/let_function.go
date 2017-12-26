package ast

import (
	"fmt"
	"strings"

	"github.com/coel-lang/coel/src/lib/debug"
	"github.com/kr/text"
)

// DefFunction represents a let-function statement node in ASTs.
type DefFunction struct {
	name      string
	signature Signature
	lets      []interface{}
	body      interface{}
	info      debug.Info
}

// NewDefFunction creates a DefFunction from its function name, signature,
// internal let statements, and a body expression.
func NewDefFunction(name string, sig Signature, lets []interface{}, expr interface{}, i debug.Info) DefFunction {
	return DefFunction{name, sig, lets, expr, i}
}

// Name returns a name of a function defined by the let-function statement.
func (f DefFunction) Name() string {
	return f.name
}

// Signature returns a signature of a function defined by the let-function statement.
func (f DefFunction) Signature() Signature {
	return f.signature
}

// Lets returns let statements contained in the let-function statement.
// Returned values should be LetVar or DefFunction.
func (f DefFunction) Lets() []interface{} {
	return f.lets
}

// Body returns a body expression of a function defined by the let-function statement.
func (f DefFunction) Body() interface{} {
	return f.body
}

// DebugInfo returns debug information of a function defined by the let-function statement.
func (f DefFunction) DebugInfo() debug.Info {
	return f.info
}

func (f DefFunction) String() string {
	ss := make([]string, 0, len(f.lets))

	for _, l := range f.lets {
		ss = append(ss, text.Indent(l.(fmt.Stringer).String(), "\t"))
	}

	return fmt.Sprintf("(def (%v %v)\n%v\t%v)\n", f.name, f.signature, strings.Join(ss, ""), f.body)
}
