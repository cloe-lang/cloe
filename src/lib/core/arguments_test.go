package core

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewArguments(t *testing.T) {
	NewArguments([]PositionalArgument{
		NewPositionalArgument(Nil, false),
		NewPositionalArgument(EmptyList, true),
		NewPositionalArgument(Nil, false),
		NewPositionalArgument(EmptyList, true),
	}, nil, nil)
}

func TestArgumentsEmpty(t *testing.T) {
	for _, a := range []Arguments{
		NewArguments(nil, nil, []Value{Nil}),
		NewArguments(nil, nil, []Value{NewDictionary([]KeyValue{{Nil, Nil}})}),
	} {
		v := a.empty()
		assert.NotNil(t, v)
		v = EvalPure(v)
		t.Logf("%#v\n", v)
		_, ok := v.(ErrorType)
		assert.True(t, ok)
	}
}

func TestArgumentsMerge(t *testing.T) {
	a := NewArguments([]PositionalArgument{NewPositionalArgument(NewList(Nil), true)}, nil, nil)
	a = a.Merge(a)
	assert.Equal(t, NewNumber(2), EvalPure(PApp(Size, a.restPositionals())))
}

func TestArgumentsRestKeywords(t *testing.T) {
	a := NewArguments(nil, nil, []Value{NewError("MyError", "")})
	assert.Equal(t, "MyError", EvalPure(a.restKeywords()).(ErrorType).Name())
}
