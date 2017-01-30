package vm

import (
	"github.com/stretchr/testify/assert"
	"testing"
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

func TestDictionarySet(t *testing.T) {
	for _, k := range []*Thunk{
		True, False, Nil, NewNumber(42), NewString("risp"),
	} {
		_, ok := App(Set, EmptyDictionary, k, Nil).Eval().(dictionaryType)
		assert.True(t, ok)
	}
}

func TestDictionaryGet(t *testing.T) {
	for _, kvs := range kvss {
		d := EmptyDictionary

		for i, kv := range kvs {
			t.Logf("Setting a %vth key...\n", i)
			d = App(Set, d, kv[0], kv[1])
		}

		assert.Equal(t, len(kvs), dictionarySize(d))

		for i, kv := range kvs {
			t.Logf("Getting a %vth value...\n", i)

			k, v := kv[0], kv[1]

			t.Log(k.Eval())

			if e, ok := App(Get, d, k).Eval().(errorType); ok {
				t.Log(e.message)
			}

			assert.True(t, testEqual(App(Get, d, k), v))
		}
	}
}

func TestDictionaryToList(t *testing.T) {
	for i, kvs := range kvss {
		t.Log("TestDictionaryToList START", i)
		d := EmptyDictionary

		for i, kv := range kvs {
			t.Logf("Setting a %vth key...\n", i)
			d = App(Set, d, kv[0], kv[1])
		}

		assert.Equal(t, len(kvs), dictionarySize(d))

		l := App(ToList, d)

		for i := 0; i < len(kvs); i, l = i+1, App(Rest, l) {
			kv := App(First, l)
			k := App(First, kv)
			lv := App(First, App(Rest, kv))
			dv := App(Get, d, k)

			t.Log("Key:", k.Eval())
			t.Log("LIST Value:", lv.Eval())
			t.Log("DICT Value:", dv.Eval())

			assert.True(t, testEqual(lv, dv))
		}

		assert.Equal(t, l.Eval().(listType), emptyList)
	}
}

func TestDictionaryWithDuplicateKeys(t *testing.T) {
	ks := []*Thunk{
		True, False, Nil, NewNumber(0), NewNumber(1), NewNumber(42),
		NewNumber(2049), NewString("runner"), NewString("lisp"),
	}

	d := EmptyDictionary

	for _, i := range []int{0, 1, 2, 2, 7, 3, 0, 4, 6, 1, 1, 4, 5, 6, 0, 2, 8, 8} {
		d = App(Set, d, ks[i], ks[i])
	}

	assert.Equal(t, len(ks), dictionarySize(d))

	for _, k := range ks {
		assert.True(t, testEqual(App(Get, d, k), k))
	}
}

func dictionarySize(d *Thunk) int {
	return int(d.Eval().(dictionaryType).Size())
}
