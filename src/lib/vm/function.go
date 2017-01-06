package vm

type Function func(List) *Thunk

func (f Function) Call(l List) *Thunk {
	return f(l)
}
