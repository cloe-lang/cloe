package core

import (
	"fmt"
	"strings"

	"github.com/coel-lang/coel/src/lib/debug"
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

// Catch returns a dictionary containing a name and message of a catched error,
// or nil otherwise.
var Catch = NewLazyFunction(
	NewSignature([]string{"error"}, nil, "", nil, nil, ""),
	func(ts ...*Thunk) Value {
		v := ts[0].Eval()
		err, ok := v.(ErrorType)

		if !ok {
			return Nil
		}

		return NewDictionary([]KeyValue{
			{NewString("name"), NewString(err.name)},
			{NewString("message"), NewString(err.message)},
		})
	})

// Chain chains 2 errors with debug information.
func (e ErrorType) Chain(i debug.Info) ErrorType {
	return ErrorType{e.name, e.message, append(e.callTrace, i)}
}

// Name returns a name of an error.
func (e ErrorType) Name() string {
	return e.name
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

// Error is implemented for error built-in interface.
func (e ErrorType) Error() string {
	return e.Lines()
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
func TypeError(v Value, typ string) *Thunk {
	s, err := StrictDump(v)

	if err != nil {
		return err
	}

	return NewError("TypeError", "%s is not a %s.", s, typ)
}

// NotBoolError creates an error value for an invalid value which is not a
// bool.
func NotBoolError(v Value) *Thunk {
	return TypeError(v, "bool")
}

// NotDictionaryError creates an error value for an invalid value which is not
// a dictionary.
func NotDictionaryError(v Value) *Thunk {
	return TypeError(v, "dictionary")
}

// NotListError creates an error value for an invalid value which is not a
// list.
func NotListError(v Value) *Thunk {
	return TypeError(v, "list")
}

// NotNumberError creates an error value for an invalid value which is not a
// number.
func NotNumberError(v Value) *Thunk {
	return TypeError(v, "number")
}

// NotIntError creates an error value for a number value which is not an
// integer.
func NotIntError(n NumberType) *Thunk {
	return TypeError(n, "integer")
}

// NotStringError creates an error value for an invalid value which is not a
// string.
func NotStringError(v Value) *Thunk {
	return TypeError(v, "string")
}

// NotCallableError creates an error value for an invalid value which is not a
// callable.
func NotCallableError(v Value) *Thunk {
	return TypeError(v, "callable")
}

// NotCollectionError creates an error value for an invalid value which is not
// a collection.
func NotCollectionError(v Value) *Thunk {
	return TypeError(v, "collection")
}

// NotEffectError creates an error value for a pure value which is expected to be an effect value.
func NotEffectError(v Value) *Thunk {
	return TypeError(v, "effect")
}

// ImpureFunctionError creates an error value for execution of an impure function.
func ImpureFunctionError(v Value) *Thunk {
	return TypeError(v, "pure value")
}

// OutOfRangeError creates an error value for an out-of-range index to a list.
func OutOfRangeError() *Thunk {
	return NewError("OutOfRangeError", "Index is out of range.")
}

func notComparableError(v Value) *Thunk {
	return TypeError(v, "comparable")
}

// NotOrderedError creates an error value for an invalid value which is not ordered.
func NotOrderedError(v Value) *Thunk {
	return TypeError(v, "ordered")
}

func emptyListError() *Thunk {
	return ValueError("The list is empty. You cannot apply rest.")
}

func keyNotFoundError(v Value) *Thunk {
	s, err := StrictDump(v)

	if err != nil {
		return err
	}

	return NewError("KeyNotFoundError", "The key %s is not found in a dictionary.", s)
}

// Error creates an error value with an error name and message.
var Error = NewLazyFunction(
	NewSignature([]string{"name", "messasge"}, nil, "", nil, nil, ""),
	func(ts ...*Thunk) Value {
		v := ts[0].Eval()
		n, ok := v.(StringType)

		if !ok {
			return NotStringError(v)
		}

		v = ts[1].Eval()
		m, ok := v.(StringType)

		if !ok {
			return NotStringError(v)
		}

		return ErrorType{string(n), string(m), []debug.Info{debug.NewGoInfo(1)}}
	})
