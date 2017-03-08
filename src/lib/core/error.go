package core

import (
	"fmt"
	"github.com/raviqqe/tisp/src/lib/debug"
)

type ErrorType struct {
	name, message string
	callTrace     []debug.Info
}

func NewError(n, m string, xs ...interface{}) *Thunk {
	return Normal(ErrorType{
		name:    n,
		message: fmt.Sprintf(m, xs...),
	})
}

func TypeError(o Object, typ string) *Thunk {
	n, m := "TypeError", "%#v is not %s"

	if e, ok := o.(ErrorType); ok {
		return Normal(e)
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

func NotDictionaryError(o Object) *Thunk {
	return TypeError(o, "Dictionary")
}

func NotStringError(o Object) *Thunk {
	return TypeError(o, "String")
}
