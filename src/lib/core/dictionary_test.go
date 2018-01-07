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
		{NewString("go"), NewList(EmptyList, Nil, NewNumber(123))},
		{False, NewNumber(42)},
		{True, False},
		{NewNumber(2), NewString("Mr. Value")},
	},
}

func TestDictionaryInsert(t *testing.T) {
	for _, k := range []*Thunk{
		True, False, Nil, NewNumber(42), NewString("coel"),
	} {
		_, ok := PApp(Insert, EmptyDictionary, k, Nil).Eval().(DictionaryType)
		assert.True(t, ok)
	}
}

func TestDictionaryInsertFail(t *testing.T) {
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

func TestDictionaryIndexFail(t *testing.T) {
	for _, th := range []*Thunk{
		PApp(EmptyDictionary, Nil),
		PApp(PApp(Insert, EmptyDictionary, NewList(OutOfRangeError()), Nil), NewList(Nil)),
	} {
		v := th.Eval()
		t.Logf("%#v", v)
		_, ok := v.(ErrorType)
		assert.True(t, ok)
	}
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
	assert.Zero(t, d.Size())
}

func TestDictionaryDeleteFail(t *testing.T) {
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

		assert.True(t, l.Eval().(ListType).Empty())
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
		{`{"foo" "bar"}`, NewDictionary([]KeyValue{{NewString("foo"), NewString("bar")}})},
	} {
		assert.Equal(t, StringType(xs.expected), PApp(ToString, xs.thunk).Eval())
	}
}

func TestDictionaryStringFail(t *testing.T) {
	for _, th := range []*Thunk{
		NewDictionary([]KeyValue{{Nil, OutOfRangeError()}}),
		NewDictionary([]KeyValue{{Nil, NewList(OutOfRangeError())}}),
	} {
		v := PApp(ToString, th).Eval()
		t.Logf("%#v", v)
		_, ok := v.(ErrorType)
		assert.True(t, ok)
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

func TestDictionaryInclude(t *testing.T) {
	for _, test := range []struct {
		dictionary *Thunk
		key        *Thunk
		answer     BoolType
	}{
		{EmptyDictionary, Nil, false},
		{PApp(Insert, EmptyDictionary, False, Nil), False, true},
		{PApp(Insert, PApp(Insert, EmptyDictionary, NewNumber(42), Nil), False, Nil), NewNumber(42), true},
		{PApp(Insert, PApp(Insert, EmptyDictionary, NewNumber(42), Nil), False, Nil), NewNumber(2049), false},
	} {
		assert.Equal(t, test.answer, PApp(Include, test.dictionary, test.key).Eval().(BoolType))
	}
}

func TestDictionaryMerge(t *testing.T) {
	d1 := EmptyDictionary
	d2kvs := make([][2]*Thunk, 0)

	for _, kvs := range kvss {
		d := EmptyDictionary

		for _, kv := range kvs {
			d = PApp(Insert, d, kv[0], kv[1])
		}

		d1 = PApp(Merge, d1, d)
		d2kvs = append(d2kvs, kvs...)
	}

	d2 := EmptyDictionary

	for _, kv := range d2kvs {
		d2 = PApp(Insert, d2, kv[0], kv[1])
	}

	assert.True(t, testEqual(d1, d2))
}

func TestDictionaryError(t *testing.T) {
	for _, th := range []*Thunk{
		PApp(
			NewDictionary([]KeyValue{{OutOfRangeError(), Nil}}),
			Nil),
		PApp(
			NewDictionary([]KeyValue{{Nil, Nil}}),
			OutOfRangeError()),
		PApp(
			Insert,
			NewDictionary([]KeyValue{{OutOfRangeError(), Nil}}),
			Nil),
		PApp(
			Insert,
			NewDictionary([]KeyValue{{Nil, Nil}}),
			OutOfRangeError()),
		PApp(
			Merge,
			EmptyDictionary,
			NewDictionary([]KeyValue{{OutOfRangeError(), Nil}})),
		PApp(
			Include,
			NewDictionary([]KeyValue{{OutOfRangeError(), Nil}}),
			Nil),
		PApp(
			ToString,
			NewDictionary([]KeyValue{{OutOfRangeError(), Nil}})),
		PApp(
			ToString,
			NewDictionary([]KeyValue{{Nil, OutOfRangeError()}})),
		PApp(
			Delete,
			NewDictionary([]KeyValue{{OutOfRangeError(), Nil}}),
			Nil),
	} {
		v := th.Eval()
		t.Log(v)
		_, ok := v.(ErrorType)
		assert.True(t, ok)
	}
}
