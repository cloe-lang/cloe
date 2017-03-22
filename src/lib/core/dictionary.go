package core

import (
	"strings"

	"github.com/raviqqe/tisp/src/lib/rbt"
	"github.com/raviqqe/tisp/src/lib/util"
)

// DictionaryType represents a dictionary in the language.
type DictionaryType struct {
	rbt.Dictionary
}

// EmptyDictionary is a thunk of an empty dictionary.
var EmptyDictionary = NewDictionary(nil, []*Thunk{})

// NewDictionary creates a dictionary from keys of values and their
// corresponding values of thunks.
func NewDictionary(ks []Value, vs []*Thunk) *Thunk {
	if len(ks) != len(vs) {
		util.Fail("Number of keys doesn't match with number of values.")
	}

	d := Normal(DictionaryType{rbt.NewDictionary(less)})

	for i, k := range ks {
		d = PApp(Insert, d, Normal(k), vs[i])
	}

	return d
}

func (d DictionaryType) insert(ts ...*Thunk) (result Value) {
	defer func() {
		if r := recover(); r != nil {
			result = r
		}
	}()

	if len(ts) != 2 {
		return NumArgsError("insert", "3 if a collection is a dictionary")
	}

	v := ts[0].Eval()

	if _, ok := v.(ordered); !ok {
		return notOrderedError(v)
	}

	return d.Insert(v, ts[1])
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

	k, ok := v.(ordered)

	if !ok {
		return notOrderedError(v)
	}

	if v, ok := d.Search(k); ok {
		return v
	}

	return NewError(
		"KeyNotFoundError",
		"The key %v is not found in a dictionary.", k)
}

func notOrderedError(k Value) *Thunk {
	return TypeError(k, "ordered")
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

func (d DictionaryType) equal(e equalable) Value {
	return d.toList().(ListType).equal(e.(DictionaryType).toList().(ListType))
}

func (d DictionaryType) toList() Value {
	k, v, rest := d.FirstRest()

	if k == nil {
		return emptyList
	}

	return cons(
		NewList(Normal(k), v),
		PApp(ToList, Normal(rest)))
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

func (d DictionaryType) size() Value {
	return NumberType(d.Size())
}

func (d DictionaryType) include(v Value) Value {
	_, ok := d.Search(v)
	return NewBool(ok)
}
