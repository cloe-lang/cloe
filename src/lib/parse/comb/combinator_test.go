package comb

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNotChar(t *testing.T) {
	s := NewState("a")
	x, err := s.NotChar(' ')()
	assert.Equal(t, 'a', x)
	assert.Equal(t, nil, err)
}

func TestNotCharFail(t *testing.T) {
	s := NewState(" ")
	x, err := s.NotChar(' ')()
	assert.Equal(t, nil, x)
	assert.NotEqual(t, nil, err)
}

func TestMany(t *testing.T) {
	for _, str := range []string{"", "  "} {
		s := NewState(str)
		x, err := s.Many(s.Char(' '))()

		t.Logf("%#v", x)

		assert.NotEqual(t, nil, x)
		assert.Equal(t, nil, err)
	}
}

func TestManyFail(t *testing.T) {
	for _, str := range []string{"="} {
		s := NewState(str)
		x, err := s.Exhaust(s.Many(func() (interface{}, error) {
			x, err := s.String("=")()

			if err != nil {
				return nil, err
			}

			if x.(string) == "=" {
				return nil, fmt.Errorf("Invalid word")
			}

			return x, nil
		}))()

		t.Logf("%#v", x)

		assert.Equal(t, nil, x)
		assert.NotEqual(t, nil, err)
	}
}

func testMany1Space(str string) (interface{}, error) {
	s := NewState(str)
	return s.Many1(s.Char(' '))()
}

func TestMany1(t *testing.T) {
	x, err := testMany1Space(" ")

	t.Logf("%#v", x)

	assert.NotEqual(t, nil, x)
	assert.Equal(t, nil, err)
}

func TestXFailMany1(t *testing.T) {
	x, err := testMany1Space("")

	t.Log(err)

	assert.Equal(t, nil, x)
	assert.NotEqual(t, nil, err)
}

func TestMany1Nest(t *testing.T) {
	s := NewState("    ")
	x, err := s.Many1(s.Many1(s.Char(' ')))()

	t.Logf("%#v", x)

	assert.NotEqual(t, nil, x)
	assert.Equal(t, nil, err)
}

func testOr(str string) (interface{}, error) {
	s := NewState(str)
	return s.Or(s.Char('a'), s.Char('b'))()
}

func TestOr(t *testing.T) {
	for _, str := range []string{"a", "b"} {
		x, err := testOr(str)

		t.Logf("%#v", x)

		assert.NotEqual(t, nil, x)
		assert.Equal(t, nil, err)
	}
}

func TestXFailOr(t *testing.T) {
	x, err := testOr("c")

	t.Log(err)

	assert.Equal(t, nil, x)
	assert.NotEqual(t, nil, err)
}

func TestMaybeSuccess(t *testing.T) {
	s := NewState("foo")
	x, err := s.Maybe(s.String("foo"))()

	t.Log(x)

	assert.Equal(t, "foo", x)
	assert.Equal(t, nil, err)
}

func TestMaybeFailure(t *testing.T) {
	s := NewState("bar")
	x, err := s.Maybe(s.String("foo"))()

	t.Log(x)

	assert.Equal(t, nil, x)
	assert.Equal(t, nil, err)
}

func TestExhaustWithErroneousParser(t *testing.T) {
	s := NewState("")
	_, err := s.Exhaust(s.String("foo"))()
	assert.NotEqual(t, nil, err)
}

func TestStringify(t *testing.T) {
	str := "foo"
	s := NewState(str)
	x, err := s.Exhaust(s.Stringify(s.And(s.String(str))))()
	assert.Equal(t, str, x)
	assert.Equal(t, nil, err)
}

func TestStringifyFail(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fail()
		}
	}()

	stringify(42)
}

func TestVoid(t *testing.T) {
	s := NewState("foo")
	x, err := s.Void(s.String("foo"))()
	assert.Equal(t, nil, x)
	assert.Equal(t, nil, err)
}
