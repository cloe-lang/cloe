package core

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsOrdered(t *testing.T) {
	for _, th := range []*Thunk{
		NewNumber(42),
		NewString("foo"),
		NewList(Nil, True),
	} {
		assert.True(t, bool(PApp(IsOrdered, th).Eval().(BoolType)))
	}

	for _, th := range []*Thunk{
		Nil,
		True,
		False,
		NewDictionary(nil, nil),
		emptyListError(),
	} {
		assert.True(t, !bool(PApp(IsOrdered, th).Eval().(BoolType)))
	}
}
