package core

import (
	"testing"

	"github.com/stretchr/testify/assert"
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

func TestBoolToString(t *testing.T) {
	test := func(s string, b bool) {
		assert.Equal(t, StringType(s), PApp(ToString, NewBool(b)).Eval())
	}

	test("true", true)
	test("false", false)
}
