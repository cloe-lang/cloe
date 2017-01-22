package ir

type LetConst struct {
	name string
	expr Expression
}

type LetFunction struct {
	name string
	body Expression
}

type Output struct {
	expr Expression
}

// For LetConst and Output,
// Expression = string | []Expression
// For LetFunction,
// Expression = int | string | []Expression
// int is used as a reference to ith argument of a function.
type Expression interface{}
