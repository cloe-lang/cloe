package vm

import "github.com/mediocregopher/seq"

type Dictionary struct{ hashMap *seq.HashMap }

func NewDictionary(ks []Object, vs []*Thunk) *Thunk {
	if len(ks) == len(vs) {
		return NewError("Number of keys doesn't match with number of values.")
	}

	d := Dictionary{seq.NewHashMap()}

	for i, k := range ks {
		d.Set(k, vs[i])
	}

	return Normal(d)
}

func (d Dictionary) Set(k, v interface{}) Dictionary {
	h, _ := d.hashMap.Set(k, v)
	return Dictionary{h}
}

func (d1 Dictionary) Equal(o Object) *Thunk {
	return NewBool(d1.hashMap.Equal(o.(Dictionary).hashMap))
}
