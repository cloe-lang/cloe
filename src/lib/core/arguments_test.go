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
	}, nil)
}

func TestArgumentsEmpty(t *testing.T) {
	for _, a := range []Arguments{
		NewPositionalArguments(Nil),
		NewPositionalArguments(Nil, Nil),
	} {
		v := a.checkEmptyness()
		assert.NotNil(t, v)
		v = EvalPure(v)
		t.Logf("%#v\n", v)
		_, ok := v.(*ErrorType)
		assert.True(t, ok)
	}
}

func TestArgumentsMerge(t *testing.T) {
	a := NewArguments([]PositionalArgument{NewPositionalArgument(NewList(Nil), true)}, nil)
	a = a.Merge(a)
	assert.Equal(t, NewNumber(2), EvalPure(PApp(Size, a.restPositionals())))
}

func TestArgumentsRestKeywords(t *testing.T) {
	a := NewArguments(nil, []KeywordArgument{NewKeywordArgument("", NewError("MyError", ""))})
	assert.Equal(t, "MyError", EvalPure(a.restKeywords()).(*ErrorType).Name())
}
