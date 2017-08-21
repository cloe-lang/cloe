package comb

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestChars(t *testing.T) {
	s := NewState("b")
	x, err := s.Chars("abc")()
	assert.Equal(t, 'b', x)
	assert.Equal(t, nil, err)
}

func TestCharsFail(t *testing.T) {
	s := NewState("d")
	x, err := s.Chars("abc")()
	assert.Equal(t, nil, x)
	assert.NotEqual(t, nil, err)
}

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

func TestWrap(t *testing.T) {
	s := NewState("abc")
	x, err := s.Wrap(s.String("a"), s.String("b"), s.String("c"))()
	assert.Equal(t, "b", x)
	assert.Equal(t, nil, err)
}

func TestPrefix(t *testing.T) {
	s := NewState("abc")
	x, err := s.Prefix(s.String("ab"), s.String("c"))()
	assert.Equal(t, "c", x)
	assert.Equal(t, nil, err)
}

func TestPrefixFail(t *testing.T) {
	s := NewState("abc")
	x, err := s.Prefix(s.String("ad"), s.String("c"))()
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
		}), exhaustError)()

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

func TestMany1Fail(t *testing.T) {
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

func TestOrFail(t *testing.T) {
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
	_, err := s.Exhaust(s.String("foo"), exhaustError)()
	assert.NotEqual(t, nil, err)
}

func TestStringify(t *testing.T) {
	str := "foo"
	s := NewState(str)
	x, err := s.Exhaust(s.Stringify(s.And(s.String(str))), exhaustError)()
	assert.Equal(t, str, x)
	assert.Equal(t, nil, err)
}

func TestStringifyFail(t *testing.T) {
	defer func() {
		assert.NotEqual(t, nil, recover())
	}()

	stringify(42)
}

func TestLazy(t *testing.T) {
	s := NewState("foo")
	x, err := s.Lazy(func() Parser { return s.String("foo") })()
	assert.Equal(t, "foo", x)
	assert.Equal(t, nil, err)
}

func TestVoid(t *testing.T) {
	s := NewState("foo")
	x, err := s.Void(s.String("foo"))()
	assert.Equal(t, nil, x)
	assert.Equal(t, nil, err)
}

func exhaustError(State) error {
	return fmt.Errorf("Parsing error")
}
