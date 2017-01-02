package types

type Callable interface {
	Call(Dictionary) Object
}
