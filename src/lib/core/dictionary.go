package core

import (
	"strings"

	"github.com/raviqqe/hamt"
)

// DictionaryType represents a dictionary in the language.
type DictionaryType struct {
	hamt.Map
}

// Eval evaluates a value into a WHNF.
func (d *DictionaryType) eval() Value {
	return d
}

var (
	emtpyDictionary = DictionaryType{hamt.NewMap()}

	// EmptyDictionary is a thunk of an empty dictionary.
	EmptyDictionary = &emtpyDictionary
)

// KeyValue is a pair of a key and value inserted into dictionaries.
type KeyValue struct {
	Key, Value Value
}

// NewDictionary creates a dictionary from keys of values and their
// corresponding values of thunks.
func NewDictionary(kvs []KeyValue) Value {
	d := Value(EmptyDictionary)

	for _, kv := range kvs {
		d = PApp(Assign, d, kv.Key, kv.Value)
	}

	return d
}

func (d *DictionaryType) assign(k Value, v Value) Value {
	e, err := evalEntry(k)

	if err != nil {
		return err
	}

	return &DictionaryType{d.Map.Insert(e, v)}
}

func (d *DictionaryType) index(k Value) Value {
	v, err := d.find(k)

	if err != nil {
		return err
	}

	return v
}

func (d *DictionaryType) find(k Value) (Value, Value) {
	e, err := evalEntry(k)

	if err != nil {
		return nil, err
	}

	v := d.Map.Find(e)

	if v == nil {
		return nil, keyNotFoundError(k)
	}

	return v.(Value), nil
}

func (d *DictionaryType) toList() Value {
	k, v, rest := d.FirstRest()

	if k == nil {
		return EmptyList
	}

	return cons(
		NewList(k.(Value), v.(Value)),
		PApp(ToList, &DictionaryType{rest}))
}

func (d *DictionaryType) merge(vs ...Value) Value {
	for _, v := range vs {
		dd, err := EvalDictionary(v)

		if err != nil {
			return err
		}

		d = &DictionaryType{d.Merge(dd.Map)}
	}

	return d
}

func (d *DictionaryType) delete(k Value) Value {
	e, err := evalEntry(k)

	if err != nil {
		return err
	}

	return &DictionaryType{d.Map.Delete(e)}
}

func (d *DictionaryType) compare(c comparable) int {
	return compare(d.toList(), c.(*DictionaryType).toList())
}

func (d *DictionaryType) string() Value {
	ss := []string{}

	for d.Size() != 0 {
		k, v, m := d.FirstRest()
		d = &DictionaryType{m}

		sk, err := StrictDump(k.(Value))

		if err != nil {
			return err
		}

		sv, err := StrictDump(EvalPure(v.(Value)))

		if err != nil {
			return err
		}

		ss = append(ss, string(sk), string(sv))
	}

	return NewString("{" + strings.Join(ss, " ") + "}")
}

func (d *DictionaryType) size() Value {
	return NewNumber(float64(d.Size()))
}

func (d *DictionaryType) include(k Value) Value {
	e, err := evalEntry(k)

	if err != nil {
		return err
	}

	return NewBoolean(d.Include(e))
}

func evalEntry(v Value) (hamt.Entry, Value) {
	e, ok := v.(hamt.Entry)

	if !ok {
		return nil, TypeError(v, "hashable")
	}

	return e, nil
}
