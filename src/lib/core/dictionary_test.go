package core

import (
	"math/rand"
	"testing"

	"github.com/stretchr/testify/assert"
)

var kvss = [][][2]*Thunk{
	{{True, False}},
	{{Nil, NewNumber(42)}},
	{{False, NewNumber(42)}, {True, NewNumber(13)}},
	{
		{False, NewNumber(42)},
		{True, False},
		{NewNumber(2), NewString("Mr. Value")},
	},
	{
		{NewString("go"), NewList(NewList(), Nil, NewNumber(123))},
		{False, NewNumber(42)},
		{True, False},
		{NewNumber(2), NewString("Mr. Value")},
	},
}

func TestDictionaryInsert(t *testing.T) {
	for _, k := range []*Thunk{
		True, False, Nil, NewNumber(42), NewString("tisp"),
	} {
		_, ok := PApp(Insert, EmptyDictionary, k, Nil).Eval().(DictionaryType)
		assert.True(t, ok)
	}
}

func TestXFailDictionaryInsert(t *testing.T) {
	l := NewList(NewError("you", "failed."))
	v := PApp(Insert, PApp(Insert, EmptyDictionary, l, Nil), l, Nil).Eval()
	_, ok := v.(ErrorType)
	t.Logf("%#v", v)
	assert.True(t, ok)
}

func TestDictionaryIndex(t *testing.T) {
	for _, kvs := range kvss {
		d := EmptyDictionary

		for i, kv := range kvs {
			t.Logf("Insertting a %vth key...\n", i)
			d = PApp(Insert, d, kv[0], kv[1])
		}

		assert.Equal(t, len(kvs), dictionarySize(d))

		for i, kv := range kvs {
			t.Logf("Getting a %vth value...\n", i)

			k, v := kv[0], kv[1]

			t.Log(k.Eval())

			if e, ok := PApp(d, k).Eval().(ErrorType); ok {
				t.Log(e.Lines())
			}

			assert.True(t, testEqual(PApp(d, k), v))
		}
	}
}

func TestXFailDictionaryIndex(t *testing.T) {
	l := NewList(NewError("you", "failed."))
	v := PApp(PApp(Insert, EmptyDictionary, l, Nil), l, Nil).Eval()
	_, ok := v.(ErrorType)
	t.Logf("%#v", v)
	assert.True(t, ok)
}

func TestDictionaryDeletable(t *testing.T) {
	t.Log(collection(EmptyDictionary.Eval().(collection)))
}

func TestDictionaryDelete(t *testing.T) {
	k := NewNumber(42)
	v := PApp(Delete, PApp(Insert, EmptyDictionary, k, Nil), k).Eval()
	d, ok := v.(DictionaryType)
	t.Logf("%#v", v)
	assert.True(t, ok)
	assert.Equal(t, 0, d.Size())
}

func TestXFailDictionaryDelete(t *testing.T) {
	l1 := NewList(NewError("you", "failed."))
	l2 := NewList(NewNumber(42))
	v := PApp(Delete, PApp(Insert, EmptyDictionary, l1, Nil), l2).Eval()
	_, ok := v.(ErrorType)
	t.Logf("%#v", v)
	assert.True(t, ok)
}

func TestDictionaryToList(t *testing.T) {
	for i, kvs := range kvss {
		t.Log("TestDictionaryToList START", i)
		d := EmptyDictionary

		for i, kv := range kvs {
			t.Logf("Insertting a %vth key...\n", i)
			d = PApp(Insert, d, kv[0], kv[1])
		}

		assert.Equal(t, len(kvs), dictionarySize(d))

		l := PApp(ToList, d)

		for i := 0; i < len(kvs); i, l = i+1, PApp(Rest, l) {
			kv := PApp(First, l)
			k := PApp(First, kv)
			lv := PApp(First, PApp(Rest, kv))
			dv := PApp(d, k)

			t.Log("Key:", k.Eval())
			t.Log("LIST Value:", lv.Eval())
			t.Log("DICT Value:", dv.Eval())

			assert.True(t, testEqual(lv, dv))
		}

		assert.Equal(t, l.Eval().(ListType), emptyList)
	}
}

func TestDictionaryWithDuplicateKeys(t *testing.T) {
	ks := []*Thunk{
		True, False, Nil, NewNumber(0), NewNumber(1), NewNumber(42),
		NewNumber(2049), NewString("runner"), NewString("lisp"),
	}

	d := EmptyDictionary

	for _, i := range []int{0, 1, 2, 2, 7, 3, 0, 4, 6, 1, 1, 4, 5, 6, 0, 2, 8, 8} {
		d = PApp(Insert, d, ks[i], ks[i])
	}

	assert.Equal(t, len(ks), dictionarySize(d))

	for _, k := range ks {
		assert.True(t, testEqual(PApp(d, k), k))
	}
}

func dictionarySize(d *Thunk) int {
	return int(d.Eval().(DictionaryType).Size())
}

func TestDictionaryEqual(t *testing.T) {
	kvs := [][2]*Thunk{
		{True, Nil},
		{False, NewList(NewNumber(123))},
		{Nil, NewList(NewNumber(123), NewNumber(456))},
		{NewNumber(42), NewString("foo")},
	}

	ds := []*Thunk{EmptyDictionary, EmptyDictionary}

	for i := range ds {
		for _, j := range rand.Perm(len(kvs)) {
			ds[i] = PApp(Insert, ds[i], kvs[j][0], kvs[j][1])
		}
	}

	assert.Equal(t, 4, ds[0].Eval().(DictionaryType).Size())
	assert.True(t, testEqual(ds[0], ds[1]))
}

func TestDictionaryLess(t *testing.T) {
	kvs := [][2]*Thunk{
		{True, Nil},
		{False, NewList(NewNumber(123))},
	}

	ds := []*Thunk{EmptyDictionary, EmptyDictionary}

	for i := range ds {
		for _, j := range rand.Perm(len(kvs)) {
			ds[i] = PApp(Insert, ds[i], kvs[j][0], kvs[j][1])
		}
	}

	ds[1] = PApp(Insert, ds[1], Nil, Nil)

	assert.Equal(t, 2, ds[0].Eval().(DictionaryType).Size())
	assert.Equal(t, 3, ds[1].Eval().(DictionaryType).Size())
	assert.True(t, testLess(ds[0], ds[1]))
}

func TestDictionaryToString(t *testing.T) {
	for _, xs := range []struct {
		expected string
		thunk    *Thunk
	}{
		{"{}", EmptyDictionary},
		{"{true nil}", PApp(Insert, EmptyDictionary, True, Nil)},
		{"{false nil true nil}", PApp(Insert, PApp(Insert, EmptyDictionary, True, Nil), False, Nil)},
	} {
		assert.Equal(t, StringType(xs.expected), PApp(ToString, xs.thunk).Eval())
	}
}

func TestDictionarySize(t *testing.T) {
	for _, test := range []struct {
		dictionary *Thunk
		size       NumberType
	}{
		{EmptyDictionary, 0},
		{PApp(Insert, EmptyDictionary, True, Nil), 1},
		{PApp(Insert, PApp(Insert, EmptyDictionary, True, Nil), False, Nil), 2},
	} {
		assert.Equal(t, test.size, PApp(Size, test.dictionary).Eval().(NumberType))
	}
}
