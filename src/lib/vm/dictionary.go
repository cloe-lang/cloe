package vm

import "github.com/mediocregopher/seq"

type Dictionary struct{ hashMap *seq.HashMap }

func NewDictionary(ks []Object, vs []*Thunk) Dictionary {
	if len(ks) != len(vs) {
		panic("Number of keys doesn't match with number of values.")
	}

	d := Dictionary{seq.NewHashMap()}

	for i, k := range ks {
		d.Set(k, vs[i])
	}

	return d
}

func (d Dictionary) Set(k, v Object) Object {
	h, _ := d.hashMap.Set((interface{})(k), (interface{})(v))
	return Dictionary{h}
}

func (d1 Dictionary) Equal(e Equalable) Object {
	// TODO: Use ToList and compare them as Lists
	return NewBool(d1.hashMap.Equal(e.(Dictionary).hashMap))
}
