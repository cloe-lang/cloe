package vm

import "github.com/mediocregopher/seq"

type nilType struct{}

var Nil = Normal(nilType{})

func (n nilType) equal(e equalable) Object {
	return True
}

// seq.Setable

func (n nilType) Hash(i uint32) uint32 {
	return i % seq.ARITY
}

func (n1 nilType) Equal(o interface{}) bool {
	_, ok := o.(nilType)
	return ok
}
