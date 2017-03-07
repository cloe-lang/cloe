package comb

import "fmt"

type Error string

func NewError(s string, xs ...interface{}) error {
	return fmt.Errorf(s, xs...)
}
