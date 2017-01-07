package vm

import "fmt"

type Error string

func NewError(s string, xs ...interface{}) *Thunk {
	return Normal(Error(fmt.Sprintf(s, xs...)))
}

func NotCallableError(o Object) *Thunk {
	return TypeError(o, "Callable")
}

func TypeError(o Object, typ string) *Thunk {
	return NewError("%#v is not %s", o, typ)
}

func NumArgsError(f, condition string) *Thunk {
	return NewError("Number of arguments to %s must be %s.", f, condition)
}

func ChainError(e *Thunk, s string, xs ...interface{}) *Thunk {
	// TODO: Error { name string, messsage string, child *Thunk }
	return nil
}

func isError(o Object) bool {
	_, ok := o.(Error)
	return ok
}
