package vm

type Nil struct{}

func NewNil() Nil {
	return Nil{}
}

func NilThunk() *Thunk {
	return Normal(NewNil())
}

func (n Nil) Equal(e Equalable) Object {
	return True
}
