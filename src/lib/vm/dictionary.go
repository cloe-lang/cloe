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
		d.set(k, vs[i])
	}

	return Normal(d)
}

var Set = NewLazyFunction(func(ts ...*Thunk) Object {
	if len(ts) != 3 {
		return NumArgsError("set", "3")
	}

	o := ts[0].Eval()
	d, ok := o.(dictionaryType)

	if !ok {
		return notDictionaryError(o)
	}

	return d.set(ts[1].Eval(), ts[2])
})

func (d dictionaryType) set(k Object, v *Thunk) Object {
	if _, ok := k.(seq.Setable); !ok {
		return notSetableError(k)
	}

	h, _ := d.hashMap.Set(k, v)
	return dictionaryType{h}
}

var Get = NewLazyFunction(func(ts ...*Thunk) Object {
	if len(ts) != 2 {
		return NumArgsError("get", "2")
	}

	o := ts[0].Eval()
	d, ok := o.(dictionaryType)

	if !ok {
		return notDictionaryError(o)
	}

	return d.get(ts[1].Eval())
})

func (d dictionaryType) get(k Object) *Thunk {
	if _, ok := k.(seq.Setable); !ok {
		return notSetableError(k)
	}

	if v, ok := d.hashMap.Get(k); ok {
		return v.(*Thunk)
	}

	return NewError(
		"KeyNotFoundError",
		"The key %v is not found in a dictionary.", k)
}

func notDictionaryError(o Object) *Thunk {
	return TypeError(o, "Dictionary")
}

func notSetableError(k Object) *Thunk {
	return TypeError(k, "Setable")
}

func (d dictionaryType) equal(e equalable) Object {
	// TODO: Use ToList and compare them as Lists
	return rawBool(d.hashMap.Equal(e.(dictionaryType).hashMap))
}

func (d dictionaryType) toList() Object {
	kv, rest, ok := d.hashMap.FirstRestKV()

	if !ok {
		return emptyList
	}

	return cons(
		NewList(Normal(kv.Key.(Object)), kv.Val.(*Thunk)),
		App(ToList, Normal(dictionaryType{rest})))
}
