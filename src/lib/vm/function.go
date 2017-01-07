package vm

type Function func(...*Thunk) *Thunk

func (f Function) Call(ts ...*Thunk) *Thunk {
	return f(ts...)
}
