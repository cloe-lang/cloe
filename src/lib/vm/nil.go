package vm

type nilType struct{}

var Nil = Normal(nilType{})

func (n nilType) equal(e equalable) Object {
	return True
}
