package vm

type Callable interface {
	Call(List) *Thunk
}

type Equalable interface {
	Equal(t1, t2 *Thunk) *Thunk
}
