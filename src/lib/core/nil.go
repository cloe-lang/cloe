package core

// NilType represents a nil value. You know.
type NilType struct{}

// Eval evaluates a value into a WHNF.
func (n NilType) eval() Value {
	return n
}

// Nil is the evil or million-dollar mistake.
var Nil = NilType{}

func (NilType) compare(comparable) int {
	return 0
}

func (NilType) string() Value {
	return NewString("nil")
}
