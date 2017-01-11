package vm

import "github.com/mediocregopher/seq"

type dictionaryType struct{ hashMap *seq.HashMap }

func NewDictionary(ks []Object, vs []*Thunk) *Thunk {
	if len(ks) != len(vs) {
		panic("Number of keys doesn't match with number of values.")
	}

	d := dictionaryType{seq.NewHashMap()}

	for i, k := range ks {
		d.Set(k, vs[i])
	}

	return Normal(d)
}

func (d dictionaryType) Set(k, v Object) Object {
	h, _ := d.hashMap.Set((interface{})(k), (interface{})(v))
	return dictionaryType{h}
}

func (d1 dictionaryType) Equal(e Equalable) Object {
	// TODO: Use ToList and compare them as Lists
	return rawBool(d1.hashMap.Equal(e.(dictionaryType).hashMap))
}
