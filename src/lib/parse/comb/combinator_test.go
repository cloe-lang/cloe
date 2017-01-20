package comb

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMany(t *testing.T) {
	for _, str := range []string{"", "  "} {
		s := NewState(str)
		result, err := s.Many(s.Char(' '))()

		t.Logf("%#v", result)

		assert.NotEqual(t, result, nil)
		assert.Equal(t, err, nil)
	}
}

func testMany1Space(str string) (interface{}, error) {
	s := NewState(str)
	return s.Many1(s.Char(' '))()
}

func TestMany1(t *testing.T) {
	result, err := testMany1Space(" ")

	t.Logf("%#v", result)

	assert.NotEqual(t, result, nil)
	assert.Equal(t, err, nil)
}

func TestMany1Fail(t *testing.T) {
	result, err := testMany1Space("")

	t.Log(err)

	assert.Equal(t, result, nil)
	assert.NotEqual(t, err, nil)
}

func TestMany1Nest(t *testing.T) {
	s := NewState("    ")
	result, err := s.Many1(s.Many1(s.Char(' ')))()

	t.Logf("%#v", result)

	assert.NotEqual(t, result, nil)
	assert.Equal(t, err, nil)
}

func testOr(str string) (interface{}, error) {
	s := NewState(str)
	return s.Or(s.Char('a'), s.Char('b'))()
}

func TestOr(t *testing.T) {
	for _, str := range []string{"a", "b"} {
		result, err := testOr(str)

		t.Logf("%#v", result)

		assert.NotEqual(t, result, nil)
		assert.Equal(t, err, nil)
	}
}

func TestOrFail(t *testing.T) {
	result, err := testOr("c")

	t.Log(err)

	assert.Equal(t, result, nil)
	assert.NotEqual(t, err, nil)
}
