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
		{PApp(NewList(NewList(Nil)), NewNumber(0), NewNumber(0)), Nil},
		{
			PApp(
				NewDictionary(
					[]Value{Nil.Eval()},
					[]*Thunk{NewDictionary([]Value{True.Eval()}, []*Thunk{False})}),
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
				NewPositionalArgument(PApp(Prepend, Nil, NewError("FooError", "Hi!")), true),
			},
			nil,
			nil)).Eval().(ErrorType)

	assert.True(t, ok)
	assert.Equal(t, "FooError", e.Name())
}
