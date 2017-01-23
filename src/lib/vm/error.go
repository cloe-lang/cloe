package vm

import "fmt"

type errorType struct {
	name, message string
	child         *errorType
}

func NewError(n, m string, xs ...interface{}) *Thunk {
	return chainedError(nil, n, m, xs...)
}

func chainedError(e *errorType, n, m string, xs ...interface{}) *Thunk {
	return Normal(errorType{
		name:    n,
		message: fmt.Sprintf(m, xs...),
		child:   e,
	})
}

func TypeError(o Object, typ string) *Thunk {
	n, m := "TypeError", "%#v is not %s"

	if e, ok := o.(errorType); ok {
		return chainedError(&e, n, m, o, typ)
	}

	return NewError(n, m, o, typ)
}

func NotCallableError(o Object) *Thunk {
	return TypeError(o, "Callable")
}

func NumArgsError(f, condition string) *Thunk {
	return NewError(
		"NumArgsError",
		"Number of arguments to %s must be %s.", f, condition)
}

func ValueError(m string) *Thunk {
	return NewError("ValueError", m)
}
