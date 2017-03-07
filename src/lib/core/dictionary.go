package core

import (
	"github.com/raviqqe/tisp/src/lib/rbt"
	"github.com/raviqqe/tisp/src/lib/util"
	"strings"
)

type DictionaryType struct{ rbt.Dictionary }

var EmptyDictionary = NewDictionary([]Object{}, []*Thunk{})

func NewDictionary(ks []Object, vs []*Thunk) *Thunk {
	if len(ks) != len(vs) {
		util.Fail("Number of keys doesn't match with number of values.")
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
	func(ts ...*Thunk) (result Object) {
		defer func() {
			if r := recover(); r != nil {
				result = r
			}
		}()

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
	func(ts ...*Thunk) (result Object) {
		defer func() {
			if r := recover(); r != nil {
				result = r
			}
		}()

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

func (d DictionaryType) delete(o Object) (result deletable, err Object) {
	defer func() {
		if r := recover(); r != nil {
			result, err = nil, r
		}
	}()

	return DictionaryType{d.Remove(o)}, nil
}

func (d DictionaryType) less(o ordered) bool {
	return less(d.toList(), o.(DictionaryType).toList())
}

func (d DictionaryType) string() Object {
	o := PApp(ToList, Normal(d)).Eval()

	if err, ok := o.(ErrorType); ok {
		return err
	}

	os, err := o.(ListType).ToObjects()

	if err != nil {
		return err.Eval()
	}

	ss := make([]string, 2*len(os))

	for i, o := range os {
		if err, ok := o.(ErrorType); ok {
			return err
		}

		os, err := o.(ListType).ToObjects()

		if err != nil {
			return err
		}

		for j, o := range os {
			if err, ok := o.(ErrorType); ok {
				return err
			}

			o = PApp(ToString, Normal(o)).Eval()

			if err, ok := o.(ErrorType); ok {
				return err
			}

			ss[2*i+j] = string(o.(StringType))
		}
	}

	return StringType("{" + strings.Join(ss, " ") + "}")
}
