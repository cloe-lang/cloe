package ast

import "github.com/raviqqe/tisp/src/lib/debug"

// MutualRecursion represents a definition of mutually-recursive functions.
type MutualRecursion struct {
	letFunctions []LetFunction
	info         debug.Info
}

// NewMutualRecursion creates a mutual recursion node from mutually-recursive
// functions.
func NewMutualRecursion(fs []LetFunction, i debug.Info) MutualRecursion {
	return MutualRecursion{fs, i}
}

// LetFunctions returns let-function statements in a mutual recursion
// definition.
func (mr MutualRecursion) LetFunctions() []LetFunction {
	return mr.letFunctions
}

// DebugInfo returns debug information of mutually-recursive function definition.
func (mr MutualRecursion) DebugInfo() debug.Info {
	return mr.info
}
