package vm

type Nil struct{}

var nill = Normal(Nil{})

func NewNil() *Thunk {
	return nill
}

func (n Nil) Equal(e Equalable) *Thunk {
	return True
}
