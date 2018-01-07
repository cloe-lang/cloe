package core

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCollectionIncludeWithErrorElement(t *testing.T) {
	v := PApp(Include, NewList(Nil), OutOfRangeError()).Eval()
	_, ok := v.(ErrorType)
	assert.True(t, ok)
}

func TestCollectionFunctionsError(t *testing.T) {
	for _, th := range []*Thunk{
		PApp(Index, Nil, Nil),
		PApp(Include, Nil, Nil),
		PApp(Merge, Nil),
		PApp(Size, Nil),
		PApp(ToList, Nil),
	} {
		v := th.Eval()
		err, ok := v.(ErrorType)
		assert.True(t, ok)
		assert.Equal(t, "TypeError", err.name)
	}
}

func TestIndexChain(t *testing.T) {
	for _, ths := range [][2]*Thunk{
		{PApp(NewList(NewList(Nil)), NewNumber(1), NewNumber(1)), Nil},
		{
			PApp(
				NewDictionary([]KeyValue{{Nil, NewDictionary([]KeyValue{{True, False}})}}),
				Nil,
				True),
			False,
		},
	} {
		assert.Equal(t, ths[1].Eval(), ths[0].Eval())
	}
}

func TestIndexWithInvalidRestArguments(t *testing.T) {
	e, ok := App(
		NewList(Nil),
		NewArguments(
			[]PositionalArgument{NewPositionalArgument(NewError("FooError", "Hi!"), true)},
			nil,
			nil)).Eval().(ErrorType)

	assert.True(t, ok)
	assert.Equal(t, "FooError", e.Name())

	e, ok = App(
		NewList(Nil),
		NewArguments(
			[]PositionalArgument{
				NewPositionalArgument(Nil, false),
				NewPositionalArgument(NewError("FooError", "Hi!"), true),
			},
			nil,
			nil)).Eval().(ErrorType)

	assert.True(t, ok)
	assert.Equal(t, "FooError", e.Name())
}

func TestInsertChain(t *testing.T) {
	for _, ts := range [][2]*Thunk{
		{NewList(True, False), PApp(Insert, EmptyList, NewNumber(1), True, NewNumber(2), False)},
		{NewList(True, False), PApp(Insert, EmptyList, NewNumber(1), False, NewNumber(1), True)},
	} {
		assert.True(t, testEqual(ts[0], ts[1]))
	}
}
