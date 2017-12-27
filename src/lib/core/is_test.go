package core

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsOrdered(t *testing.T) {
	for _, th := range []*Thunk{
		NewNumber(42),
		NewString("foo"),
		EmptyList,
		NewList(NewNumber(42)),
		NewList(NewNumber(42), EmptyList),
		NewList(NewNumber(42), EmptyList, NewList(NewNumber(42), NewString("foo"))),
	} {
		assert.True(t, bool(PApp(IsOrdered, th).Eval().(BoolType)))
	}

	for _, th := range []*Thunk{
		Nil,
		True,
		False,
		EmptyDictionary,
		NewList(Nil, True),
		NewList(NewNumber(42), Nil),
		NewList(NewNumber(42), EmptyList, NewList(Nil)),
		NewList(NewNumber(42), EmptyList, NewList(NewNumber(42), Nil, NewString("foo"))),
	} {
		assert.False(t, bool(PApp(IsOrdered, th).Eval().(BoolType)))
	}
}

func TestIsOrderedError(t *testing.T) {
	for _, th := range []*Thunk{
		emptyListError(),
		PApp(Add, Nil),
		NewList(emptyListError()),
		PApp(Prepend, NewNumber(42), emptyListError()),
	} {
		_, ok := PApp(IsOrdered, th).Eval().(ErrorType)
		assert.True(t, ok)
	}
}
