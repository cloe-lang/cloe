package core

import (
	"strings"

	"github.com/coel-lang/coel/src/lib/rbt"
)

// DictionaryType represents a dictionary in the language.
type DictionaryType struct {
	rbt.Dictionary
}

// Eval evaluates a value into a WHNF.
func (d DictionaryType) eval() Value {
	return d
}

// EmptyDictionary is a thunk of an empty dictionary.
var EmptyDictionary = DictionaryType{rbt.NewDictionary(compare)}

// KeyValue is a pair of a key and value inserted into dictionaries.
type KeyValue struct {
	Key, Value Value
}

// NewDictionary creates a dictionary from keys of values and their
// corresponding values of thunks.
func NewDictionary(kvs []KeyValue) Value {
	d := Value(EmptyDictionary)

	for _, kv := range kvs {
		d = PApp(Insert, d, kv.Key, kv.Value)
	}

	return d
}

func (d DictionaryType) call(args Arguments) Value {
	return Index.call(NewPositionalArguments(d).Merge(args))
}

func (d DictionaryType) index(v Value) (result Value) {
	defer func() {
		if r := recover(); r != nil {
			result = r.(Value)
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

func (d DictionaryType) insert(k Value, v Value) (result Value) {
	defer func() {
		if r := recover(); r != nil {
			result = r.(Value)
		}
	}()

	if _, ok := k.(comparable); !ok {
		return notComparableError(k)
	}

	return d.Insert(k, v)
}

func (d DictionaryType) toList() (result Value) {
	defer func() {
		if r := recover(); r != nil {
			result = r.(Value)
		}
	}()

	k, v, rest := d.FirstRest()

	if k == nil {
		return EmptyList
	}

	return cons(
		NewList(k, v),
		PApp(ToList, rest))
}

func (d DictionaryType) merge(vs ...Value) (result Value) {
	defer func() {
		if r := recover(); r != nil {
			result = r.(Value)
		}
	}()

	for _, v := range vs {
		dd, err := EvalDictionary(v)

		if err != nil {
			return err
		}

		d = d.Merge(dd)
	}

	return d
}

func (d DictionaryType) delete(v Value) (result Value) {
	defer func() {
		if r := recover(); r != nil {
			result = r.(Value)
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
			result = r.(Value)
		}
	}()

	ss := []string{}

	for !d.Empty() {
		var k, v Value
		k, v, d = d.FirstRest()

		sk, err := StrictDump(k)

		if err != nil {
			return err
		}

		sv, err := StrictDump(EvalPure(v))

		if err != nil {
			return err
		}

		ss = append(ss, string(sk), string(sv))
	}

	return NewString("{" + strings.Join(ss, " ") + "}")
}

func (d DictionaryType) size() Value {
	return NewNumber(float64(d.Size()))
}

func (d DictionaryType) include(v Value) (result Value) {
	defer func() {
		if r := recover(); r != nil {
			result = r.(Value)
		}
	}()

	_, ok := d.Search(v)
	return NewBool(ok)
}

// Insert wraps rbt.Dictionary.Insert().
func (d DictionaryType) Insert(k Value, v Value) DictionaryType {
	return DictionaryType{d.Dictionary.Insert(k, v)}
}

// Search wraps rbt.Dictionary.Search().
func (d DictionaryType) Search(k Value) (Value, bool) {
	v, ok := d.Dictionary.Search(k)

	if !ok {
		return nil, false
	}

	return v.(Value), true
}

// Remove wraps rbt.Dictionary.Remove().
func (d DictionaryType) Remove(k Value) DictionaryType {
	return DictionaryType{d.Dictionary.Remove(k)}
}

// FirstRest wraps rbt.Dictionary.FirstRest().
func (d DictionaryType) FirstRest() (Value, Value, DictionaryType) {
	k, v, rest := d.Dictionary.FirstRest()
	d = DictionaryType{rest}

	if k == nil {
		return nil, nil, d
	}

	return k.(Value), v.(Value), d
}

// Merge wraps rbt.Dictionary.Merge().
func (d DictionaryType) Merge(dd DictionaryType) DictionaryType {
	return DictionaryType{d.Dictionary.Merge(dd.Dictionary)}
}
