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

func TestDictionarySetMethod(t *testing.T) {
	for _, th := range []*Thunk{
		True, False, Nil, NewNumber(42), NewString("risp"),
	} {
		_, ok := EmptyDictionary.Eval().(dictionaryType).Set(th.Eval(), Nil).(dictionaryType)
		assert.True(t, ok)
	}
}

func TestDictionaryGetMethod(t *testing.T) {
	for _, kvs := range kvss {
		d := EmptyDictionary.Eval().(dictionaryType)

		for i, kv := range kvs {
			t.Logf("Setting a %vth key...\n", i)
			d = d.Set(extractKeyAndValue(kv)).(dictionaryType)
		}

		assert.Equal(t, len(kvs), int(d.hashMap.Size()))

		for i, kv := range kvs {
			t.Logf("Getting a %vth value...\n", i)

			k, v := extractKeyAndValue(kv)

			t.Log(k)

			if e, ok := d.Get(k).Eval().(errorType); ok {
				t.Log(e.message.Eval())
			}

			assert.True(t, testEqual(d.Get(k), v))
		}
	}
}

func TestDictionaryToList(t *testing.T) {
	for i, kvs := range kvss {
		t.Log("TestDictionaryToList START", i)
		d := EmptyDictionary.Eval().(dictionaryType)

		for i, kv := range kvs {
			t.Logf("Setting a %vth key...\n", i)
			d = d.Set(extractKeyAndValue(kv)).(dictionaryType)
		}

		assert.Equal(t, len(kvs), int(d.hashMap.Size()))

		l := App(ToList, Normal(d))

		for i := 0; i < len(kvs); i, l = i+1, App(Rest, l) {
			kv := App(First, l)
			k := App(First, kv).Eval()
			lv := App(First, App(Rest, kv))
			dv := d.Get(k)

			t.Log("Key:", k)
			t.Log("LIST Value:", lv.Eval())
			t.Log("DICT Value:", dv.Eval())

			assert.True(t, testEqual(lv, dv))
		}

		assert.Equal(t, l.Eval().(listType), emptyList)
	}
}

func extractKeyAndValue(kv [2]*Thunk) (Object, *Thunk) {
	return kv[0].Eval(), kv[1]
}
