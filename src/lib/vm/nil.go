package vm

func Nil() *Thunk {
	return Normal(nil)
}
