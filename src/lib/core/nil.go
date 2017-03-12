package core

// NilType is a type of Nil. You know.
type NilType struct{}

// Nil is the evil or million-dollar mistake.
var Nil = Normal(NilType{})

func (n NilType) equal(e equalable) Value {
	return True
}

func (NilType) less(o ordered) bool {
	return false
}

func (NilType) string() Value {
	return StringType("nil")
}
