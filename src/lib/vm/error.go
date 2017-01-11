package vm

import "fmt"

type errorType struct {
	// TODO: child can be *errorType
	name, message, child *Thunk
}

func internalError(n, m string, xs ...interface{}) *Thunk {
	return chainedError(NilThunk(), n, m, xs...)
}

func chainedError(e *Thunk, n, m string, xs ...interface{}) *Thunk {
	return Normal(errorType{
		name:    Normal(NewString(n)),
		message: Normal(NewString(fmt.Sprintf(m, xs...))),
		child:   e,
	})
}

func TypeError(o Object, typ string) *Thunk {
	n, m := "TypeError", "%#v is not %s"

	if e, ok := o.(errorType); ok {
		return chainedError(Normal(e), n, m, o, typ)
	}

	return internalError(n, m, o, typ)
}

func NotCallableError(o Object) *Thunk {
	return TypeError(o, "Callable")
}

func NumArgsError(f, condition string) *Thunk {
	return internalError(
		"NumArgsError",
		"Number of arguments to %s must be %s.", f, condition)
}

func ValueError(m string) *Thunk {
	return internalError("ValueError", m)
}
