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
