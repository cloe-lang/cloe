package vm

type nilType struct{}

var Nil = Normal(nilType{})

func (n nilType) Equal(e Equalable) Object {
	return True
}
