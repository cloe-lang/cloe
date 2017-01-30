package vm

import "../rbt"

type dictionaryType struct{ rbt.Dictionary }

var EmptyDictionary = NewDictionary([]Object{}, []*Thunk{})

func NewDictionary(ks []Object, vs []*Thunk) *Thunk {
	if len(ks) != len(vs) {
		panic("Number of keys doesn't match with number of values.")
	}

	d := Normal(dictionaryType{rbt.NewDictionary(less)})

	for i, k := range ks {
		d = App(Set, d, Normal(k), vs[i])
	}

	return d
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

	k := ts[1].Eval()

	if _, ok := k.(ordered); !ok {
		return notOrderedError(k)
	}

	return dictionaryType{d.Insert(k, ts[2])}
})

var Get = NewLazyFunction(func(ts ...*Thunk) Object {
	if len(ts) != 2 {
		return NumArgsError("get", "2")
	}

	o := ts[0].Eval()
	d, ok := o.(dictionaryType)

	if !ok {
		return notDictionaryError(o)
	}

	o = ts[1].Eval()
	k, ok := o.(ordered)

	if !ok {
		return notOrderedError(o)
	}

	if v, ok := d.Search(k); ok {
		return v.(*Thunk)
	}

	return NewError(
		"KeyNotFoundError",
		"The key %v is not found in a dictionary.", k)
})

func notDictionaryError(o Object) *Thunk {
	return TypeError(o, "Dictionary")
}

func notOrderedError(k Object) *Thunk {
	return TypeError(k, "Ordered")
}

func (d dictionaryType) equal(e equalable) Object {
	return d.toList().(listType).equal(e.(dictionaryType).toList().(listType))
}

func (d dictionaryType) toList() Object {
	k, v, rest := d.FirstRest()

	if k == nil {
		return emptyList
	}

	return cons(
		NewList(Normal(k.(Object)), v.(*Thunk)),
		App(ToList, Normal(dictionaryType{rest})))
}
