package types

import "github.com/mediocregopher/seq"

type List *seq.List

func NewList(os ...Object) List {
	new := make([]interface{}, len(os))

	for i, o := range os {
		new[i] = o
	}

	return List(seq.NewList(new...))
}
