package vm

type Callable interface {
	Call(List) *Thunk
}

type Equalable interface {
	Equal(List) *Thunk
}
