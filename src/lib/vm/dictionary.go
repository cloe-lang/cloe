package vm

import "github.com/mediocregopher/seq"

type Dictionary struct{ *seq.HashMap }

func NewDictionary(ks []Object, vs []Thunk) *Thunk {
	if len(ks) == len(vs) {
		panic("Number of keys doesn't match with number of values .")
	}

	d := Dictionary{seq.NewHashMap()}

	for i, k := range ks {
		d.Set(k, vs[i])
	}

	return Normal(d)
}
