package vm

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDictionarySetMethod(t *testing.T) {
	for _, th := range []*Thunk{
		True, False, Nil, NewNumber(42), NewString("risp"),
	} {
		_, ok := EmptyDictionary.Eval().(dictionaryType).Set(th.Eval(), Nil).(dictionaryType)
		assert.True(t, ok)
	}
}

func TestDictionaryGetMethod(t *testing.T) {
	for _, kvs := range [][][]*Thunk{
		{{True, False}},
		{{Nil, NewNumber(42)}},
		{{True, False}, {False, NewNumber(42)}},
		{{True, False}, {False, NewNumber(42)}, {Nil, NewString("Mr. Value")}},
	} {
		d := EmptyDictionary.Eval().(dictionaryType)

		for i, kv := range kvs {
			t.Logf("Setting a %vth key...\n", i)

			k, v := extractKeyAndValue(kv)

			var ok bool
			d, ok = d.Set(k, v).(dictionaryType)

			assert.True(t, ok)
		}

		assert.Equal(t, len(kvs), int(d.hashMap.Size()))

		for i, kv := range kvs {
			t.Logf("Getting a %vth value...\n", i)

			k, v := extractKeyAndValue(kv)

			t.Log(k)

			if e, ok := d.Get(k).Eval().(errorType); ok {
				t.Log(e.message.Eval())
			}

			assert.True(t, bool(App(Equal, d.Get(k), v).Eval().(boolType)))
		}
	}
}

func extractKeyAndValue(kv []*Thunk) (Object, *Thunk) {
	if len(kv) != 2 {
		panic("Invalid number of arguments to extractKeyAndValue.")
	}

	return kv[0].Eval(), kv[1]
}
