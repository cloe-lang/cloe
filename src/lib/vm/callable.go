package vm

type Callable interface {
	Call(List) *Thunk
}
