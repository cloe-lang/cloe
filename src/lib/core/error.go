package core

import (
	"fmt"
	"github.com/raviqqe/tisp/src/lib/debug"
	"strings"
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

func (e ErrorType) Lines() string {
	ss := make([]string, 0, len(e.callTrace))

	for i := range e.callTrace {
		ss = append(ss, e.callTrace[len(e.callTrace)-1-i].Lines())
	}

	return strings.Join(ss, "") + e.name + ": " + e.message + "\n"
}

func TypeError(o Object, typ string) *Thunk {
	if e, ok := o.(ErrorType); ok {
		return Normal(e)
	}

	return NewError("TypeError", "%#v is not a %s.", o, typ)
}

func NumArgsError(f, condition string) *Thunk {
	return NewError("NumArgsError", "Number of arguments to %s must be %s.", f, condition)
}

func ValueError(m string) *Thunk {
	return NewError("ValueError", m)
}

func NotBoolError(o Object) *Thunk {
	return TypeError(o, "bool")
}

func NotDictionaryError(o Object) *Thunk {
	return TypeError(o, "dictionary")
}

func NotListError(o Object) *Thunk {
	return TypeError(o, "list")
}

func NotNumberError(o Object) *Thunk {
	return TypeError(o, "number")
}

func NotStringError(o Object) *Thunk {
	return TypeError(o, "string")
}

func NotCallableError(o Object) *Thunk {
	return TypeError(o, "funtion.")
}
