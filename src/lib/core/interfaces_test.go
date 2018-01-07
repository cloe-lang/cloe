package core

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func testEqual(t1, t2 *Thunk) bool {
	return compare(t1.Eval(), t2.Eval()) == 0
}

func testLess(t1, t2 *Thunk) bool {
	return t1.Eval().(comparable).compare(t2.Eval().(comparable)) < 0
}

func TestLessFail(t *testing.T) {
	defer func() {
		assert.NotNil(t, recover())
	}()

	compare(NewNumber(42).Eval(), NewError("you", "failed.").Eval())
}

func testCompare(t1, t2 *Thunk) int {
	return int(PApp(Compare, t1, t2).Eval().(NumberType))
}

func TestCompareWithInvalidValues(t *testing.T) {
	for _, ts := range [][2]*Thunk{
		{True, False},
		{Nil, Nil},
		{NewNumber(0), False},
		{NewNumber(0), Nil},
		{True, Nil},
		{NewDictionary([]KeyValue{{Nil, Nil}}), NewDictionary([]KeyValue{{Nil, Nil}})},
		{NotNumberError(Nil), NotNumberError(Nil)},
	} {
		v := PApp(Compare, ts[0], ts[1]).Eval()

		t.Log(v)

		_, ok := v.(ErrorType)
		assert.True(t, ok)
	}
}

func TestCompareErrorMessage(t *testing.T) {
	e, ok := PApp(Compare, EmptyList, NewString("foo")).Eval().(ErrorType)

	assert.True(t, ok)
	assert.Equal(t, "[] is not a string.", e.message)
}

func TestOrdered(t *testing.T) {
	for _, th := range []*Thunk{
		NewNumber(42),
		NewString("foo"),
		EmptyList,
	} {
		o, ok := th.Eval().(ordered)
		assert.True(t, ok)
		o.ordered()
	}
}

func TestNotOrdered(t *testing.T) {
	for _, th := range []*Thunk{
		Nil,
		True,
		False,
		EmptyDictionary,
		ValueError("This is not ordered."),
		Merge,
	} {
		_, ok := th.Eval().(ordered)
		assert.False(t, ok)
	}
}

func TestEqualTrue(t *testing.T) {
	for _, ts := range [][]*Thunk{
		{},
		{True},
		{True, True},
		{Nil, Nil},
		{NewNumber(42), NewNumber(42)},
		{NewNumber(42), NewNumber(42), NewNumber(42)},
	} {
		assert.True(t, bool(PApp(Equal, ts...).Eval().(BoolType)))
	}
}

func TestEqualFalse(t *testing.T) {
	for _, ts := range [][]*Thunk{
		{True, False},
		{NewNumber(0), NewNumber(42)},
		{NewNumber(0), NewNumber(0), NewNumber(42)},
		{NewNumber(0), NewNumber(42), NewNumber(42)},
		{NewList(NewNumber(42), NewNumber(42)), many42()},
		{NewList(NewNumber(42), NewNumber(42)), many42()},
		{many42(), NewList(NewNumber(42), NewNumber(42))},
	} {
		assert.True(t, !bool(PApp(Equal, ts...).Eval().(BoolType)))
	}
}

func TestEqualFail(t *testing.T) {
	for _, a := range []Arguments{
		NewPositionalArguments(NewError("", ""), Nil),
		NewPositionalArguments(Nil, NewError("", "")),
		NewArguments(
			[]PositionalArgument{NewPositionalArgument(OutOfRangeError(), true)},
			nil,
			nil),
		NewArguments(
			[]PositionalArgument{
				NewPositionalArgument(Nil, false),
				NewPositionalArgument(OutOfRangeError(), true),
			},
			nil,
			nil),
		NewArguments(
			[]PositionalArgument{
				NewPositionalArgument(Nil, false),
				NewPositionalArgument(NewList(OutOfRangeError()), true),
			},
			nil,
			nil),
	} {
		_, ok := App(Equal, a).Eval().(ErrorType)
		assert.True(t, ok)
	}
}

func many42() *Thunk {
	return PApp(Prepend, NewNumber(42), PApp(NewLazyFunction(
		NewSignature(nil, nil, "", nil, nil, ""),
		func(...*Thunk) Value {
			return many42()
		})))
}
