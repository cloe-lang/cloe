package vm

import "fmt"

type Error string

func NewError(s string, xs ...interface{}) Error {
	return Error(fmt.Sprintf(s, xs...))
}

func (e Error) Error() string {
	return string(e)
}
