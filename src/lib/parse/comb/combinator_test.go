package comb

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
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

func TestManyFail(t *testing.T) {
	for _, str := range []string{"="} {
		s := NewState(str)
		result, err := s.Exhaust(s.Many(func() (interface{}, error) {
			x, err := s.String("=")()

			if err != nil {
				return nil, err
			}

			if x.(string) == "=" {
				return nil, fmt.Errorf("Invalid word")
			}

			return x, nil
		}))()

		t.Logf("%#v", result)

		assert.Equal(t, result, nil)
		assert.NotEqual(t, err, nil)
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

func TestXFailMany1(t *testing.T) {
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

func TestXFailOr(t *testing.T) {
	result, err := testOr("c")

	t.Log(err)

	assert.Equal(t, result, nil)
	assert.NotEqual(t, err, nil)
}

func TestMaybeSuccess(t *testing.T) {
	s := NewState("foo")
	result, err := s.Maybe(s.String("foo"))()

	t.Log(result)

	assert.Equal(t, "foo", result)
	assert.Equal(t, nil, err)
}

func TestMaybeFailure(t *testing.T) {
	s := NewState("bar")
	result, err := s.Maybe(s.String("foo"))()

	t.Log(result)

	assert.Equal(t, nil, result)
	assert.Equal(t, nil, err)
}

func TestExhaustWithErroneousParser(t *testing.T) {
	s := NewState("")
	_, err := s.Exhaust(s.String("foo"))()
	assert.NotEqual(t, err, nil)
}

func TestStringify(t *testing.T) {
	str := "foo"
	s := NewState(str)
	result, err := s.Exhaust(s.Stringify(s.And(s.String(str))))()
	assert.Equal(t, str, result.(string))
	assert.Equal(t, err, nil)
}

func TestStringifyFail(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fail()
		}
	}()

	stringify(42)
}
