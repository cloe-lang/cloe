package vm

import "fmt"

type Error string

func NewError(s string, xs ...interface{}) *Thunk {
	return Normal(Error(fmt.Sprintf(s, xs...)))
}

func NotCallableError(o Object) *Thunk {
	return NewError("%#v is not callable.", o)
}

func ChainError(e *Thunk, s string, xs ...interface{}) *Thunk {
	// TODO: Error { name string, messsage string, child *Thunk }
	return nil
}

func isError(o Object) bool {
	_, ok := o.(Error)
	return ok
}
