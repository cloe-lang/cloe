package core

import (
	"fmt"
	"strings"

	"github.com/raviqqe/tisp/src/lib/debug"
)

// ErrorType represents errors in the language and traces function calls for
// debugging.
type ErrorType struct {
	name, message string
	callTrace     []debug.Info
}

// NewError creates an error value from its name and a formatted message.
func NewError(n, m string, xs ...interface{}) *Thunk {
	return Normal(ErrorType{
		name:    n,
		message: fmt.Sprintf(m, xs...),
	})
}

// Lines returns multi-line string representation of an error which can be
// printed as is to stdout or stderr.
func (e ErrorType) Lines() string {
	ss := make([]string, 0, len(e.callTrace))

	for i := range e.callTrace {
		ss = append(ss, e.callTrace[len(e.callTrace)-1-i].Lines())
	}

	return strings.Join(ss, "") + e.name + ": " + e.message + "\n"
}

// NumArgsError creates an error value for an invalid number of arguments
// passed to a function.
func NumArgsError(f, condition string) *Thunk {
	return NewError("NumArgsError", "Number of arguments to %s must be %s.", f, condition)
}

// ValueError creates an error value for some invalid value detected at runtime.
func ValueError(m string, xs ...interface{}) *Thunk {
	return NewError("ValueError", m, xs...)
}

// TypeError creates an error value for an invalid type.
func TypeError(o Object, typ string) *Thunk {
	if e, ok := o.(ErrorType); ok {
		return Normal(e)
	}

	return NewError("TypeError", "%#v is not a %s.", o, typ)
}

// NotBoolError creates an error value for an invalid value which is not a
// bool.
func NotBoolError(o Object) *Thunk {
	return TypeError(o, "bool")
}

// NotDictionaryError creates an error value for an invalid value which is not
// a dictionary.
func NotDictionaryError(o Object) *Thunk {
	return TypeError(o, "dictionary")
}

// NotListError creates an error value for an invalid value which is not a
// list.
func NotListError(o Object) *Thunk {
	return TypeError(o, "list")
}

// NotNumberError creates an error value for an invalid value which is not a
// number.
func NotNumberError(o Object) *Thunk {
	return TypeError(o, "number")
}

// NotStringError creates an error value for an invalid value which is not a
// string.
func NotStringError(o Object) *Thunk {
	return TypeError(o, "string")
}

// NotCallableError creates an error value for an invalid value which is not a
// callable.
func NotCallableError(o Object) *Thunk {
	return TypeError(o, "function.")
}

// InputError creates a thunk which represents an input error.
func InputError(m string, xs ...interface{}) *Thunk {
	return NewError("InputError", m, xs...)
}

// OutputError creates a thunk which represents an output error.
func OutputError(m string, xs ...interface{}) *Thunk {
	return NewError("OutputError", m, xs...)
}
