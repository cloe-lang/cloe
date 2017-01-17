package vm

import "fmt"

type errorType struct {
	name, message string
	child         *errorType
}

func internalError(n, m string, xs ...interface{}) *Thunk {
	return chainedError(nil, n, m, xs...)
}

func chainedError(e *errorType, n, m string, xs ...interface{}) *Thunk {
	return Normal(errorType{
		name:    n,
		message: fmt.Sprintf(m, xs...),
		child:   e,
	})
}

func typeError(o Object, typ string) *Thunk {
	n, m := "TypeError", "%#v is not %s"

	if e, ok := o.(errorType); ok {
		return chainedError(&e, n, m, o, typ)
	}

	return internalError(n, m, o, typ)
}

func notCallableError(o Object) *Thunk {
	return typeError(o, "Callable")
}

func numArgsError(f, condition string) *Thunk {
	return internalError(
		"NumArgsError",
		"Number of arguments to %s must be %s.", f, condition)
}

func valueError(m string) *Thunk {
	return internalError("ValueError", m)
}
