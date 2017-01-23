package ir

type LetConst struct {
	name string
	expr Expression
}

func NewLetConst(s string, e Expression) LetConst {
	return LetConst{name: s, expr: e}
}

type LetFunction struct {
	name string
	body Expression
}

func NewLetFunction(s string, e Expression) LetFunction {
	return LetFunction{name: s, body: e}
}

type Output struct {
	expr Expression
}

func NewOutput(e Expression) Output {
	return Output{e}
}

// For LetConst and Output,
// Expression = string | []Expression
// For LetFunction,
// Expression = int | string | []Expression
// int is used as a reference to ith argument of a function.
type Expression interface{}
