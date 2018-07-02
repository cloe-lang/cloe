package core

import (
	"hash/fnv"

	"github.com/raviqqe/hamt"
)

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

// Hash hashes a value.
func (NilType) Hash() uint32 {
	return fnv.New32().Sum32()
}

// Equal checks equality.
func (NilType) Equal(e hamt.Entry) bool {
	_, ok := e.(NilType)
	return ok
}
