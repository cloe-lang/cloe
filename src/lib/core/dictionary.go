package core

import (
	"strings"

	"github.com/raviqqe/tisp/src/lib/rbt"
	"github.com/raviqqe/tisp/src/lib/util"
)

type DictionaryType struct{ rbt.Dictionary }

var EmptyDictionary = NewDictionary([]Value{}, []*Thunk{})

func NewDictionary(ks []Value, vs []*Thunk) *Thunk {
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
	func(ts ...*Thunk) (result Value) {
		defer func() {
			if r := recover(); r != nil {
				result = r
			}
		}()

		v := ts[0].Eval()
		d, ok := v.(DictionaryType)

		if !ok {
			return NotDictionaryError(v)
		}

		k := ts[1].Eval()

		if _, ok := k.(ordered); !ok {
			return notOrderedError(k)
		}

		return DictionaryType{d.Insert(k, ts[2])}
	})

func (d DictionaryType) call(args Arguments) Value {
	return Index.Eval().(callable).call(NewPositionalArguments(Normal(d)).Merge(args))
}

func (d DictionaryType) index(v Value) (result Value) {
	defer func() {
		if r := recover(); r != nil {
			result = r
		}
	}()

	k, ok := v.(ordered)

	if !ok {
		return notOrderedError(v)
	}

	if v, ok := d.Search(k); ok {
		return v.(*Thunk)
	}

	return NewError(
		"KeyNotFoundError",
		"The key %v is not found in a dictionary.", k)
}

func notOrderedError(k Value) *Thunk {
	return TypeError(k, "ordered")
}

func (d DictionaryType) equal(e equalable) Value {
	return d.toList().(ListType).equal(e.(DictionaryType).toList().(ListType))
}

func (d DictionaryType) toList() Value {
	k, v, rest := d.FirstRest()

	if k == nil {
		return emptyList
	}

	return cons(
		NewList(Normal(k.(Value)), v.(*Thunk)),
		PApp(ToList, Normal(DictionaryType{rest})))
}

func (d DictionaryType) merge(ts ...*Thunk) Value {
	for _, t := range ts {
		go t.Eval()
	}

	for _, t := range ts {
		v := t.Eval()
		dd, ok := v.(DictionaryType)

		if !ok {
			return NotDictionaryError(v)
		}

		d = DictionaryType{d.Merge(dd.Dictionary)}
	}

	return d
}

func (d DictionaryType) delete(v Value) (result deletable, err Value) {
	defer func() {
		if r := recover(); r != nil {
			result, err = nil, r
		}
	}()

	return DictionaryType{d.Remove(v)}, nil
}

func (d DictionaryType) less(o ordered) bool {
	return less(d.toList(), o.(DictionaryType).toList())
}

func (d DictionaryType) string() Value {
	v := PApp(ToList, Normal(d)).Eval()
	l, ok := v.(ListType)

	if !ok {
		return NotListError(v)
	}

	vs, err := l.ToValues()

	if err != nil {
		return err.Eval()
	}

	ss := make([]string, 2*len(vs))

	for i, v := range vs {
		if err, ok := v.(ErrorType); ok {
			return err
		}

		vs, err := v.(ListType).ToValues()

		if err != nil {
			return err
		}

		for j, v := range vs {
			if err, ok := v.(ErrorType); ok {
				return err
			}

			v = PApp(ToString, Normal(v)).Eval()
			s, ok := v.(StringType)

			if !ok {
				return NotStringError(v)
			}

			ss[2*i+j] = string(s)
		}
	}

	return StringType("{" + strings.Join(ss, " ") + "}")
}
