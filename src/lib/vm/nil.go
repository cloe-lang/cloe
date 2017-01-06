package vm

type Nil struct{}

func NewNil() *Thunk {
	return Normal(Nil{})
}
