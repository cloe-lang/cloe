package vm

type NilType struct{}

var Nil = Normal(NilType{})

func (n NilType) equal(e equalable) Object {
	return True
}

// ordered

func (NilType) less(o ordered) bool {
	return false
}
