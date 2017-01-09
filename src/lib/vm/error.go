package vm

import "fmt"

type Error struct {
	// TODO: child can be *Error
	name, message, child *Thunk
}

func internalError(n, m string, xs ...interface{}) Error {
	return chainedError(NilThunk(), n, m, xs...)
}

func chainedError(e *Thunk, n, m string, xs ...interface{}) Error {
	return Error{
		name:    Normal(NewString(n)),
		message: Normal(NewString(fmt.Sprintf(m, xs...))),
		child:   e,
	}
}

func TypeError(o Object, typ string) Error {
	n, m := "TypeError", "%#v is not %s"

	if e, ok := o.(Error); ok {
		return chainedError(Normal(e), n, m, o, typ)
	}

	return internalError(n, m, o, typ)
}

func NotCallableError(o Object) Error {
	return TypeError(o, "Callable")
}

func NumArgsError(f, condition string) Error {
	return internalError(
		"NumArgsError",
		"Number of arguments to %s must be %s.", f, condition)
}

func ValueError(m string) Error {
	return internalError("ValueError", m)
}
