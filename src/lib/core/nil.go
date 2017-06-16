package core

// NilType represents a nil value. You know.
type NilType struct{}

// Nil is the evil or million-dollar mistake.
var Nil = Normal(NilType{})

func (NilType) compare(ordered) int {
	return 0
}

func (NilType) string() Value {
	return StringType("nil")
}
