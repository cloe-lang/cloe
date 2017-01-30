package vm

type nilType struct{}

var Nil = Normal(nilType{})

func (n nilType) equal(e equalable) Object {
	return True
}

// ordered

func (nilType) less(o ordered) bool {
	return false
}
