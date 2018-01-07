package core

import (
	"strings"

	"github.com/coel-lang/coel/src/lib/rbt"
)

// DictionaryType represents a dictionary in the language.
type DictionaryType struct {
	rbt.Dictionary
}

var (
	emptyDictionary = DictionaryType{rbt.NewDictionary(compare)}

	// EmptyDictionary is a thunk of an empty dictionary.
	EmptyDictionary = Normal(emptyDictionary)
)

// KeyValue is a pair of a key and value inserted into dictionaries.
type KeyValue struct {
	Key, Value *Thunk
}

// NewDictionary creates a dictionary from keys of values and their
// corresponding values of thunks.
func NewDictionary(kvs []KeyValue) *Thunk {
	d := EmptyDictionary

	for _, kv := range kvs {
		d = PApp(Insert, d, kv.Key, kv.Value)
	}

	return d
}

func (d DictionaryType) call(args Arguments) Value {
	return Index.Eval().(callable).call(NewPositionalArguments(Normal(d)).Merge(args))
}

func (d DictionaryType) index(v Value) (result Value) {
	defer func() {
		if r := recover(); r != nil {
			result = r
		}
	}()

	k, ok := v.(comparable)

	if !ok {
		return notComparableError(v)
	}

	if v, ok := d.Search(k); ok {
		return v
	}

	return keyNotFoundError(k)
}

func (d DictionaryType) insert(v Value, t *Thunk) (result Value) {
	defer func() {
		if r := recover(); r != nil {
			result = r
		}
	}()

	if _, ok := v.(comparable); !ok {
		return notComparableError(v)
	}

	return d.Insert(v, t)
}

func (d DictionaryType) toList() (result Value) {
	defer func() {
		if r := recover(); r != nil {
			result = r
		}
	}()

	k, v, rest := d.FirstRest()

	if k == nil {
		return emptyList
	}

	return cons(
		NewList(Normal(k), v),
		PApp(ToList, Normal(rest)))
}

func (d DictionaryType) merge(ts ...*Thunk) (result Value) {
	defer func() {
		if r := recover(); r != nil {
			result = r
		}
	}()

	for _, t := range ts {
		v := t.Eval()
		dd, ok := v.(DictionaryType)

		if !ok {
			return NotDictionaryError(v)
		}

		d = d.Merge(dd)
	}

	return d
}

func (d DictionaryType) delete(v Value) (result Value) {
	defer func() {
		if r := recover(); r != nil {
			result = r
		}
	}()

	return d.Remove(v)
}

func (d DictionaryType) compare(c comparable) int {
	return compare(d.toList(), c.(DictionaryType).toList())
}

func (d DictionaryType) string() (result Value) {
	defer func() {
		if r := recover(); r != nil {
			result = r
		}
	}()

	ss := []string{}

	for !d.Empty() {
		var (
			k Value
			v *Thunk
		)

		k, v, d = d.FirstRest()

		sk, err := StrictDump(k)

		if err != nil {
			return err
		}

		sv, err := StrictDump(v.Eval())

		if err != nil {
			return err
		}

		ss = append(ss, string(sk), string(sv))
	}

	return StringType("{" + strings.Join(ss, " ") + "}")
}

func (d DictionaryType) size() Value {
	return NumberType(d.Size())
}

func (d DictionaryType) include(v Value) (result Value) {
	defer func() {
		if r := recover(); r != nil {
			result = r
		}
	}()

	_, ok := d.Search(v)
	return NewBool(ok)
}

// Insert wraps rbt.Dictionary.Insert().
func (d DictionaryType) Insert(k Value, v *Thunk) DictionaryType {
	return DictionaryType{d.Dictionary.Insert(k, v)}
}

// Search wraps rbt.Dictionary.Search().
func (d DictionaryType) Search(k Value) (*Thunk, bool) {
	v, ok := d.Dictionary.Search(k)

	if !ok {
		return nil, false
	}

	return v.(*Thunk), true
}

// Remove wraps rbt.Dictionary.Remove().
func (d DictionaryType) Remove(k Value) DictionaryType {
	return DictionaryType{d.Dictionary.Remove(k)}
}

// FirstRest wraps rbt.Dictionary.FirstRest().
func (d DictionaryType) FirstRest() (Value, *Thunk, DictionaryType) {
	k, v, rest := d.Dictionary.FirstRest()
	d = DictionaryType{rest}

	if k == nil {
		return nil, nil, d
	}

	return k.(Value), v.(*Thunk), d
}

// Merge wraps rbt.Dictionary.Merge().
func (d DictionaryType) Merge(dd DictionaryType) DictionaryType {
	return DictionaryType{d.Dictionary.Merge(dd.Dictionary)}
}
