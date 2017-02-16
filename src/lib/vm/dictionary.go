package vm

import "../rbt"

type DictionaryType struct{ rbt.Dictionary }

var EmptyDictionary = NewDictionary([]Object{}, []*Thunk{})

func NewDictionary(ks []Object, vs []*Thunk) *Thunk {
	if len(ks) != len(vs) {
		panic("Number of keys doesn't match with number of values.")
	}

	d := Normal(DictionaryType{rbt.NewDictionary(less)})

	for i, k := range ks {
		d = PApp(Set, d, Normal(k), vs[i])
	}

	return d
}

var Set = NewLazyFunction(
	NewSignature(
		[]string{"dict", "key", "value"}, []OptionalArgument{}, "",
		[]string{}, []OptionalArgument{}, "",
	),
	func(ts ...*Thunk) Object {
		o := ts[0].Eval()
		d, ok := o.(DictionaryType)

		if !ok {
			return NotDictionaryError(o)
		}

		k := ts[1].Eval()

		if _, ok := k.(ordered); !ok {
			return notOrderedError(k)
		}

		return DictionaryType{d.Insert(k, ts[2])}
	})

var Get = NewLazyFunction(
	NewSignature(
		[]string{"dict", "key"}, []OptionalArgument{}, "",
		[]string{}, []OptionalArgument{}, "",
	),
	func(ts ...*Thunk) Object {
		o := ts[0].Eval()
		d, ok := o.(DictionaryType)

		if !ok {
			return NotDictionaryError(o)
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

func notOrderedError(k Object) *Thunk {
	return TypeError(k, "Ordered")
}

func (d DictionaryType) equal(e equalable) Object {
	return d.toList().(ListType).equal(e.(DictionaryType).toList().(ListType))
}

func (d DictionaryType) toList() Object {
	k, v, rest := d.FirstRest()

	if k == nil {
		return emptyList
	}

	return cons(
		NewList(Normal(k.(Object)), v.(*Thunk)),
		PApp(ToList, Normal(DictionaryType{rest})))
}

func (d DictionaryType) merge(ts ...*Thunk) Object {
	for _, t := range ts {
		go t.Eval()
	}

	for _, t := range ts {
		o := t.Eval()
		dd, ok := o.(DictionaryType)

		if !ok {
			return NotDictionaryError(o)
		}

		d = DictionaryType{d.Merge(dd.Dictionary)}
	}

	return d
}

// ordered

func (d DictionaryType) less(o ordered) bool {
	return less(d.toList(), o.(DictionaryType).toList())
}
