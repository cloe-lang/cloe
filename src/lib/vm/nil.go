package vm

type nilType struct{}

func NewNil() nilType {
	return nilType{}
}

func NilThunk() *Thunk {
	return Normal(NewNil())
}

func (n nilType) Equal(e Equalable) Object {
	return True
}
