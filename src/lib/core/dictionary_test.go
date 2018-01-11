package core

import (
	"math/rand"
	"testing"

	"github.com/stretchr/testify/assert"
)

var kvss = [][][2]Value{
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
	for _, k := range []Value{
		True, False, Nil, NewNumber(42), NewString("coel"),
	} {
		_, ok := EvalPure(PApp(Insert, EmptyDictionary, k, Nil)).(*DictionaryType)
		assert.True(t, ok)
	}
}

func TestDictionaryInsertFail(t *testing.T) {
	l := NewList(NewError("you", "failed."))
	v := EvalPure(PApp(Insert, PApp(Insert, EmptyDictionary, l, Nil), l, Nil))
	_, ok := v.(ErrorType)
	t.Logf("%#v", v)
	assert.True(t, ok)
}

func TestDictionaryIndex(t *testing.T) {
	for _, kvs := range kvss {
		d := Value(EmptyDictionary)

		for i, kv := range kvs {
			t.Logf("Insertting a %vth key...\n", i)
			d = PApp(Insert, d, kv[0], kv[1])
		}

		assert.Equal(t, len(kvs), dictionarySize(d))

		for i, kv := range kvs {
			t.Logf("Getting a %vth value...\n", i)

			k, v := kv[0], kv[1]

			t.Log(EvalPure(k))

			if e, ok := EvalPure(PApp(d, k)).(ErrorType); ok {
				t.Log(e.Lines())
			}

			assert.True(t, testEqual(PApp(d, k), v))
		}
	}
}

func TestDictionaryIndexFail(t *testing.T) {
	for _, v := range []Value{
		PApp(EmptyDictionary, Nil),
		PApp(PApp(Insert, EmptyDictionary, NewList(DummyError), Nil), NewList(Nil)),
		PApp(
			PApp(Insert, EmptyDictionary, NewList(Nil, DummyError), Nil),
			NewList(Nil, DummyError)),
	} {
		v := EvalPure(v)
		t.Logf("%#v", v)
		_, ok := v.(ErrorType)
		assert.True(t, ok)
	}
}

func TestDictionaryDelete(t *testing.T) {
	k := NewNumber(42)
	v := EvalPure(PApp(Delete, PApp(Insert, EmptyDictionary, k, Nil), k))
	d, ok := v.(*DictionaryType)
	t.Logf("%#v", v)
	assert.True(t, ok)
	assert.Zero(t, d.Size())
}

func TestDictionaryDeleteFail(t *testing.T) {
	v := EvalPure(PApp(
		Delete,
		PApp(Insert, EmptyDictionary, NewList(DummyError), Nil),
		NewList(NewNumber(42))))
	_, ok := v.(ErrorType)
	t.Logf("%#v", v)
	assert.True(t, ok)
}

func TestDictionaryToList(t *testing.T) {
	for i, kvs := range kvss {
		t.Log("TestDictionaryToList START", i)
		d := Value(EmptyDictionary)

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

			t.Log("Key:", EvalPure(k))
			t.Log("LIST Value:", EvalPure(lv))
			t.Log("DICT Value:", EvalPure(dv))

			assert.True(t, testEqual(lv, dv))
		}

		assert.True(t, EvalPure(l).(*ListType).Empty())
	}
}

func TestDictionaryWithDuplicateKeys(t *testing.T) {
	ks := []Value{
		True, False, Nil, NewNumber(0), NewNumber(1), NewNumber(42),
		NewNumber(2049), NewString("runner"), NewString("lisp"),
	}

	d := Value(EmptyDictionary)

	for _, i := range []int{0, 1, 2, 2, 7, 3, 0, 4, 6, 1, 1, 4, 5, 6, 0, 2, 8, 8} {
		d = PApp(Insert, d, ks[i], ks[i])
	}

	assert.Equal(t, len(ks), dictionarySize(d))

	for _, k := range ks {
		assert.True(t, testEqual(PApp(d, k), k))
	}
}

func dictionarySize(d Value) int {
	return int(EvalPure(d).(*DictionaryType).Size())
}

func TestDictionaryEqual(t *testing.T) {
	kvs := [][2]Value{
		{True, Nil},
		{False, NewList(NewNumber(123))},
		{Nil, NewList(NewNumber(123), NewNumber(456))},
		{NewNumber(42), NewString("foo")},
	}

	ds := []Value{EmptyDictionary, EmptyDictionary}

	for i := range ds {
		for _, j := range rand.Perm(len(kvs)) {
			ds[i] = PApp(Insert, ds[i], kvs[j][0], kvs[j][1])
		}
	}

	assert.Equal(t, 4, dictionarySize(ds[0]))
	assert.True(t, testEqual(ds[0], ds[1]))
}

func TestDictionaryLess(t *testing.T) {
	kvs := [][2]Value{
		{True, Nil},
		{False, NewList(NewNumber(123))},
	}

	ds := []Value{EmptyDictionary, EmptyDictionary}

	for i := range ds {
		for _, j := range rand.Perm(len(kvs)) {
			ds[i] = PApp(Insert, ds[i], kvs[j][0], kvs[j][1])
		}
	}

	ds[1] = PApp(Insert, ds[1], Nil, Nil)

	assert.Equal(t, 2, dictionarySize(ds[0]))
	assert.Equal(t, 3, dictionarySize(ds[1]))
	assert.True(t, testLess(ds[0], ds[1]))
}

func TestDictionaryToString(t *testing.T) {
	for _, c := range []struct {
		expected string
		value    Value
	}{
		{"{}", EmptyDictionary},
		{"{true nil}", PApp(Insert, EmptyDictionary, True, Nil)},
		{"{false nil true nil}", PApp(Insert, PApp(Insert, EmptyDictionary, True, Nil), False, Nil)},
		{`{"foo" "bar"}`, NewDictionary([]KeyValue{{NewString("foo"), NewString("bar")}})},
	} {
		assert.Equal(t, StringType(c.expected), EvalPure(PApp(ToString, c.value)))
	}
}

func TestDictionaryStringFail(t *testing.T) {
	for _, v := range []Value{
		NewDictionary([]KeyValue{{Nil, OutOfRangeError()}}),
		NewDictionary([]KeyValue{{Nil, NewList(OutOfRangeError())}}),
		NewDictionary([]KeyValue{{NewList(OutOfRangeError()), Nil}}),
	} {
		v := EvalPure(PApp(ToString, v))
		t.Logf("%#v", v)
		_, ok := v.(ErrorType)
		assert.True(t, ok)
	}
}

func TestDictionarySize(t *testing.T) {
	for _, test := range []struct {
		dictionary Value
		size       NumberType
	}{
		{EmptyDictionary, 0},
		{PApp(Insert, EmptyDictionary, True, Nil), 1},
		{PApp(Insert, PApp(Insert, EmptyDictionary, True, Nil), False, Nil), 2},
	} {
		assert.Equal(t, test.size, EvalPure(PApp(Size, test.dictionary)).(NumberType))
	}
}

func TestDictionaryInclude(t *testing.T) {
	for _, c := range []struct {
		dictionary Value
		key        Value
		answer     BoolType
	}{
		{EmptyDictionary, Nil, (false)},
		{PApp(Insert, EmptyDictionary, False, Nil), False, (true)},
		{PApp(Insert, PApp(Insert, EmptyDictionary, NewNumber(42), Nil), False, Nil), NewNumber(42), (true)},
		{PApp(Insert, PApp(Insert, EmptyDictionary, NewNumber(42), Nil), False, Nil), NewNumber(2049), (false)},
	} {
		assert.Equal(t, c.answer, *EvalPure(PApp(Include, c.dictionary, c.key)).(*BoolType))
	}
}

func TestDictionaryMerge(t *testing.T) {
	d1 := Value(EmptyDictionary)
	d2kvs := make([][2]Value, 0)

	for _, kvs := range kvss {
		d := Value(EmptyDictionary)

		for _, kv := range kvs {
			d = PApp(Insert, d, kv[0], kv[1])
		}

		d1 = PApp(Merge, d1, d)
		d2kvs = append(d2kvs, kvs...)
	}

	d2 := Value(EmptyDictionary)

	for _, kv := range d2kvs {
		d2 = PApp(Insert, d2, kv[0], kv[1])
	}

	assert.True(t, testEqual(d1, d2))
}

func TestDictionaryError(t *testing.T) {
	for _, v := range []Value{
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
			Merge,
			NewDictionary([]KeyValue{{NewList(OutOfRangeError()), Nil}}),
			NewDictionary([]KeyValue{{NewList(OutOfRangeError()), Nil}})),
		PApp(
			Include,
			NewDictionary([]KeyValue{{OutOfRangeError(), Nil}}),
			Nil),
		PApp(
			Include,
			NewDictionary([]KeyValue{{NewList(Nil), Nil}}),
			NewList(OutOfRangeError())),
		PApp(
			ToString,
			NewDictionary([]KeyValue{{OutOfRangeError(), Nil}})),
		PApp(
			ToString,
			NewDictionary([]KeyValue{{Nil, OutOfRangeError()}})),
		PApp(
			ToString,
			NewDictionary([]KeyValue{{NewList(OutOfRangeError()), OutOfRangeError()}})),
		PApp(
			Delete,
			NewDictionary([]KeyValue{{OutOfRangeError(), Nil}}),
			Nil),
		PApp(
			ToList,
			NewDictionary([]KeyValue{{NewList(OutOfRangeError()), Nil}})),
	} {
		v := EvalPure(v)
		t.Log(v)
		_, ok := v.(ErrorType)
		assert.True(t, ok)
	}
}
