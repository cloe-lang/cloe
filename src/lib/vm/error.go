package vm

import "fmt"

type Error string

func NewError(s string, xs ...interface{}) *Thunk {
	return Normal(Error(fmt.Sprintf(s, xs...)))
}
