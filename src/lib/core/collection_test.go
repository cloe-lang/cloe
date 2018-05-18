package core

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCollectionIncludeWithErrorElement(t *testing.T) {
	_, ok := EvalPure(PApp(Include, NewList(Nil), DummyError)).(*ErrorType)
	assert.True(t, ok)
}

func TestCollectionFunctionsError(t *testing.T) {
	for _, v := range []Value{
		PApp(Index, Nil, Nil),
		PApp(Index, Nil, Nil),
		PApp(Include, Nil, Nil),
		App(
			Insert,
			NewArguments(
				[]PositionalArgument{
					NewPositionalArgument(EmptyDictionary, false),
					NewPositionalArgument(Nil, true),
				}, nil)),
		App(
			Insert,
			NewArguments(
				[]PositionalArgument{
					NewPositionalArgument(EmptyDictionary, false),
					NewPositionalArgument(StrictPrepend([]Value{Nil}, EmptyDictionary), true),
				}, nil)),
		App(
			Insert,
			NewArguments(
				[]PositionalArgument{
					NewPositionalArgument(EmptyDictionary, false),
					NewPositionalArgument(StrictPrepend([]Value{Nil, Nil}, EmptyDictionary), true),
				}, nil)),
		App(
			Insert,
			NewArguments(
				[]PositionalArgument{
					NewPositionalArgument(EmptyDictionary, false),
					NewPositionalArgument(StrictPrepend([]Value{Nil, Nil, Nil}, EmptyDictionary), true),
				}, nil)),
		PApp(Merge, Nil),
		App(
			Merge,
			NewArguments(
				[]PositionalArgument{
					NewPositionalArgument(EmptyList, false),
					NewPositionalArgument(Nil, true),
				}, nil)),
		App(
			Merge,
			NewArguments(
				[]PositionalArgument{
					NewPositionalArgument(EmptyList, false),
					NewPositionalArgument(StrictPrepend([]Value{Nil}, EmptyDictionary), true),
				}, nil)),
		PApp(Size, Nil),
		PApp(ToList, Nil),
	} {
		err, ok := EvalPure(v).(*ErrorType)
		assert.True(t, ok)
		assert.Equal(t, "TypeError", err.name)
	}
}

func TestIndexChain(t *testing.T) {
	for _, vs := range [][2]Value{
		{
			PApp(NewList(NewList(Nil)), NewNumber(1), NewNumber(1)),
			Nil,
		},
		{
			PApp(
				NewDictionary([]KeyValue{{Nil, NewDictionary([]KeyValue{{True, False}})}}),
				Nil,
				True),
			False,
		},
	} {
		assert.Equal(t, EvalPure(vs[1]), EvalPure(vs[0]))
	}
}

func TestIndexWithInvalidRestArguments(t *testing.T) {
	e, ok := EvalPure(App(
		NewList(Nil),
		NewArguments(
			[]PositionalArgument{NewPositionalArgument(DummyError, true)},
			nil))).(*ErrorType)

	assert.True(t, ok)
	assert.Equal(t, "DummyError", e.Name())

	e, ok = EvalPure(App(
		NewList(Nil),
		NewArguments(
			[]PositionalArgument{
				NewPositionalArgument(Nil, false),
				NewPositionalArgument(DummyError, true),
			},
			nil))).(*ErrorType)

	assert.True(t, ok)
	assert.Equal(t, "DummyError", e.Name())
}

func TestInsertChain(t *testing.T) {
	for _, vs := range [][2]Value{
		{NewList(True, False), PApp(Insert, EmptyList, NewNumber(1), True, NewNumber(2), False)},
		{NewList(True, False), PApp(Insert, EmptyList, NewNumber(1), False, NewNumber(1), True)},
	} {
		assert.True(t, testEqual(vs[0], vs[1]))
	}
}
