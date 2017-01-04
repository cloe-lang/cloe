package vm

type Callable interface {
	Call(Dictionary) Object
}
