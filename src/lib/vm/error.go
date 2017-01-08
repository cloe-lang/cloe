package vm

import "fmt"

type Error struct {
	// Errors should be lazy maybe.
	name, message string
	child         *Thunk
}

func NewError(n, m string, xs ...interface{}) *Thunk {
	return ChainError(NewNil(), n, m, xs...)
}

func ChainError(e *Thunk, n, m string, xs ...interface{}) *Thunk {
	return Normal(Error{
		name:    n,
		message: fmt.Sprintf(m, xs...),
		child:   e,
	})
}

func NotCallableError(o Object) *Thunk {
	return TypeError(o, "Callable")
}

func TypeError(o Object, typ string) *Thunk {
	n, m := "TypeError", "%#v is not %s"

	if e, ok := o.(Error); ok {
		return ChainError(Normal(e), n, m, o, typ)
	}

	return NewError(n, m, o, typ)
}

func NumArgsError(f, condition string) *Thunk {
	return NewError(
		"NumArgsError",
		"Number of arguments to %s must be %s.", f, condition)
}

func ValueError(m string) *Thunk {
	return NewError("ValueError", m)
}
