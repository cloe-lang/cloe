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
	if _, ok := k.(seq.Setable); !ok {
		return notSetableError(k)
	}

	h, _ := d.hashMap.Set(k, v)
	return dictionaryType{h}
}

func (d dictionaryType) Get(k Object) *Thunk {
	if _, ok := k.(seq.Setable); !ok {
		return notSetableError(k)
	}

	if v, ok := d.hashMap.Get(k); ok {
		return v.(*Thunk)
	}

	return internalError(
		"KeyNotFoundError",
		"The key %v is not found in a dictionary.", k)
}

func (d1 dictionaryType) equal(e equalable) Object {
	// TODO: Use ToList and compare them as Lists
	return rawBool(d1.hashMap.Equal(e.(dictionaryType).hashMap))
}

func notSetableError(k Object) *Thunk {
	return typeError(k, "Setable")
}
