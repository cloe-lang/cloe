package vm

import "github.com/mediocregopher/seq"

type dictionaryType struct{ hashMap *seq.HashMap }

var EmptyDictionary = NewDictionary([]Object{}, []*Thunk{})

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

func (d dictionaryType) Set(k Object, v *Thunk) Object {
	if s, ok := k.(seq.Setable); ok {
		h, _ := d.hashMap.Set(s, v)
		return dictionaryType{h}
	}

	return typeError(k, "Setable")
}

func (d1 dictionaryType) equal(e equalable) Object {
	// TODO: Use ToList and compare them as Lists
	return rawBool(d1.hashMap.Equal(e.(dictionaryType).hashMap))
}
