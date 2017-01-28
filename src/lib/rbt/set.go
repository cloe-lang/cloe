package rbt

type Set struct{ Tree }

func NewSet(less func(interface{}, interface{}) bool) Set {
	return Set{NewTree(less)}
}
