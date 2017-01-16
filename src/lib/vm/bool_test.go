package vm

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestBoolEqual(t *testing.T) {
	for _, ts := range [][]*Thunk{
		{True, True},
		{False, False},
	} {
		assert.True(t, testEqual(ts...))
	}

	for _, ts := range [][]*Thunk{
		{True, False},
		{False, True},
	} {
		assert.True(t, !testEqual(ts...))
	}
}
