package core

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func testEqual(t1, t2 Value) bool {
	return compare(EvalPure(t1), EvalPure(t2)) == 0
}

func testLess(t1, t2 Value) bool {
	return EvalPure(t1).(comparable).compare(EvalPure(t2).(comparable)) < 0
}

func testCompare(t1, t2 Value) int {
	return int(*EvalPure(PApp(Compare, t1, t2)).(*NumberType))
}

func TestCompareFunction(t *testing.T) {
	assert.Equal(t, -1, compare(True, NewNumber(42)))
}

func TestComaprePanic(t *testing.T) {
	assert.Panics(t, func() {
		compare(NewNumber(42), NewError("you", "failed."))
	})
}

func TestCompareWithInvalidValues(t *testing.T) {
	for _, vs := range [][2]Value{
		{True, False},
		{Nil, Nil},
		{NewNumber(0), False},
		{NewNumber(0), Nil},
		{True, Nil},
		{NewDictionary([]KeyValue{{Nil, Nil}}), NewDictionary([]KeyValue{{Nil, Nil}})},
		{NotNumberError(Nil), NotNumberError(Nil)},
	} {
		v := EvalPure(PApp(Compare, vs[0], vs[1]))

		t.Log(v)

		_, ok := v.(ErrorType)
		assert.True(t, ok)
	}
}

func TestCompareErrorMessage(t *testing.T) {
	e, ok := EvalPure(PApp(Compare, EmptyList, NewString("foo"))).(ErrorType)

	assert.True(t, ok)
	assert.Equal(t, "[] is not a string.", e.message)
}

func TestOrdered(t *testing.T) {
	for _, v := range []Value{
		NewNumber(42),
		NewString("foo"),
		EmptyList,
	} {
		o, ok := EvalPure(v).(ordered)
		assert.True(t, ok)
		o.ordered()
	}
}

func TestNotOrdered(t *testing.T) {
	for _, v := range []Value{
		Nil,
		True,
		False,
		EmptyDictionary,
		ValueError("This is not ordered."),
		Merge,
	} {
		_, ok := EvalPure(v).(ordered)
		assert.False(t, ok)
	}
}

func TestEqualTrue(t *testing.T) {
	for _, vs := range [][]Value{
		{},
		{True},
		{True, True},
		{Nil, Nil},
		{NewNumber(42), NewNumber(42)},
		{NewNumber(42), NewNumber(42), NewNumber(42)},
	} {
		assert.True(t, bool(*EvalPure(PApp(Equal, vs...)).(*BoolType)))
	}
}

func TestEqualFalse(t *testing.T) {
	for _, vs := range [][]Value{
		{True, False},
		{NewNumber(0), NewNumber(42)},
		{NewNumber(0), NewNumber(0), NewNumber(42)},
		{NewNumber(0), NewNumber(42), NewNumber(42)},
		{NewList(NewNumber(42), NewNumber(42)), many42()},
		{NewList(NewNumber(42), NewNumber(42)), many42()},
		{many42(), NewList(NewNumber(42), NewNumber(42))},
	} {
		assert.True(t, !bool(*EvalPure(PApp(Equal, vs...)).(*BoolType)))
	}
}

func TestEqualFail(t *testing.T) {
	for _, a := range []Arguments{
		NewPositionalArguments(DummyError, Nil),
		NewPositionalArguments(Nil, DummyError),
		NewArguments(
			[]PositionalArgument{NewPositionalArgument(DummyError, true)},
			nil,
			nil),
		NewArguments(
			[]PositionalArgument{
				NewPositionalArgument(Nil, false),
				NewPositionalArgument(DummyError, true),
			},
			nil,
			nil),
		NewArguments(
			[]PositionalArgument{
				NewPositionalArgument(Nil, false),
				NewPositionalArgument(NewList(DummyError), true),
			},
			nil,
			nil),
	} {
		_, ok := EvalPure(App(Equal, a)).(ErrorType)
		assert.True(t, ok)
	}
}

func TestIsOrdered(t *testing.T) {
	for _, v := range []Value{
		NewNumber(42),
		NewString("foo"),
		EmptyList,
		NewList(NewNumber(42)),
		NewList(NewNumber(42), EmptyList),
		NewList(NewNumber(42), EmptyList, NewList(NewNumber(42), NewString("foo"))),
	} {
		assert.True(t, bool(*EvalPure(PApp(IsOrdered, v)).(*BoolType)))
	}

	for _, v := range []Value{
		Nil,
		True,
		False,
		EmptyDictionary,
		NewList(Nil, True),
		NewList(NewNumber(42), Nil),
		NewList(NewNumber(42), EmptyList, NewList(Nil)),
		NewList(NewNumber(42), EmptyList, NewList(NewNumber(42), Nil, NewString("foo"))),
	} {
		assert.False(t, bool(*EvalPure(PApp(IsOrdered, v)).(*BoolType)))
	}
}

func TestIsOrderedError(t *testing.T) {
	for _, v := range []Value{
		DummyError,
		PApp(Add, Nil),
		NewList(DummyError),
		PApp(Prepend, NewNumber(42), DummyError),
	} {
		_, ok := EvalPure(PApp(IsOrdered, v)).(ErrorType)
		assert.True(t, ok)
	}
}

func many42() Value {
	return PApp(Prepend, NewNumber(42), PApp(NewLazyFunction(
		NewSignature(nil, nil, "", nil, nil, ""),
		func(...Value) Value {
			return many42()
		})))
}
