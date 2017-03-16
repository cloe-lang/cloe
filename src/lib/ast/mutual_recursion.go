package ast

// MutualRecursion represents a definition of mutually-recursive functions.
type MutualRecursion struct {
	letFunctions []LetFunction
}

// NewMutualRecursion creates a mutual recursion node from mutually-recursive
// functions.
func NewMutualRecursion(fs ...LetFunction) MutualRecursion {
	return MutualRecursion{fs}
}

// LetFunctions returns let-function statements in a mutual recursion
// definition.
func (mr MutualRecursion) LetFunctions() []LetFunction {
	return mr.letFunctions
}
