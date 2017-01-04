package vm

import "github.com/mediocregopher/seq"

type List struct{ *seq.List }

func NewList(ts ...*Thunk) List {
	new := make([]interface{}, len(ts))

	for i, t := range ts {
		new[i] = t
	}

	return List{seq.NewList(new...)}
}
