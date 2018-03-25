package core

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStringEqual(t *testing.T) {
	s := NewString("foo")

	assert.True(t, testEqual(s, s))
	assert.True(t, !testEqual(s, NewString("bar")))
}

func TestStringMerge(t *testing.T) {
	v := NewString("foo")
	assert.Equal(t, v+v, EvalPure(PApp(Merge, v, v)))
}

func TestStringMergeWithNonString(t *testing.T) {
	_, ok := EvalPure(PApp(Merge, NewString("foo"), Nil)).(ErrorType)
	assert.True(t, ok)
}

func TestStringToList(t *testing.T) {
	s := "lisp"
	l := PApp(ToList, NewString(s))

	for _, r := range s {
		assert.Equal(t, NewString(string(r)), EvalPure(PApp(First, l)))
		l = PApp(Rest, l)
	}

	assert.True(t, EvalPure(l).(*ListType).Empty())
}

func TestStringIndex(t *testing.T) {
	for _, vs := range [][2]Value{
		{NewString("axc"), NewNumber(2)},
		{NewString("ああaあxbあc"), NewNumber(5)},
	} {
		assert.Equal(t, NewString("x"), EvalPure(PApp(vs[0], vs[1:]...)))
	}
}

func TestStringIndexWithInvalidIndexNumber(t *testing.T) {
	for _, i := range []float64{-1, 1.5, 100} {
		_, ok := EvalPure(PApp(NewString("foo"), NewNumber(i))).(ErrorType)
		assert.True(t, ok)
	}
}

func TestStringDelete(t *testing.T) {
	for _, c := range []struct {
		string   StringType
		index    float64
		expected StringType
	}{
		{"a", 1, ""},
		{"abc", 2, "ac"},
		{"abc", 3, "ab"},
	} {
		assert.Equal(t, c.expected, EvalPure(PApp(Delete, c.string, NewNumber(c.index))))
	}
}

func TestStringDeleteWithInvalidIndex(t *testing.T) {
	for _, i := range []float64{-1, 100} {
		_, ok := EvalPure(PApp(Delete, NewString("foo"), NewNumber(i))).(ErrorType)
		assert.True(t, ok)
	}
}

func TestStringSize(t *testing.T) {
	for _, c := range []struct {
		string StringType
		size   NumberType
	}{
		{"", 0},
		{"a", 1},
		{"ab", 2},
		{"abc", 3},
	} {
		assert.Equal(t, c.size, *EvalPure(PApp(Size, c.string)).(*NumberType))
	}
}

func TestStringInclude(t *testing.T) {
	for _, c := range []struct {
		string    StringType
		substring StringType
		answer    BoolType
	}{
		{"", "", true},
		{"a", "", true},
		{"a", "a", true},
		{"abc", "ab", true},
		{"abcdef", "cde", true},
		{"", "a", false},
		{"ab", "ac", false},
	} {
		assert.Equal(t, c.answer, *EvalPure(PApp(Include, c.string, c.substring)).(*BoolType))
	}
}

func TestStringIncludeWithNonString(t *testing.T) {
	_, ok := EvalPure(PApp(Include, NewString("foo"), Nil)).(ErrorType)
	assert.True(t, ok)
}

func TestStringInsert(t *testing.T) {
	for _, c := range []struct {
		string   StringType
		index    NumberType
		elem     StringType
		expected StringType
	}{
		{"", 1, "", ""},
		{"", 1, "a", "a"},
		{"a", 1, "b", "ba"},
		{"a", 2, "b", "ab"},
		{"ab", 1, "c", "cab"},
		{"ab", 2, "c", "acb"},
		{"ab", 3, "c", "abc"},
	} {
		assert.True(t, testEqual(c.expected, PApp(Insert, c.string, &c.index, c.elem)))
	}
}

func TestStringInsertWithInvalidIndex(t *testing.T) {
	for _, i := range []float64{-1, 0, 5, 100} {
		_, ok := EvalPure(PApp(Insert, NewString("foo"), NewNumber(i), NewString("bar"))).(ErrorType)
		assert.True(t, ok)
	}
}

func TestStringInsertWithNonString(t *testing.T) {
	_, ok := EvalPure(PApp(Insert, NewString("foo"), NewNumber(1), Nil)).(ErrorType)
	assert.True(t, ok)
}

func TestStringCompare(t *testing.T) {
	for _, c := range []struct {
		answer      int
		left, right Value
	}{
		{0, NewString("foo"), NewString("foo")},
		{1, NewString("foo"), NewString("bar")},
		{-1, NewString("bar"), NewString("foo")},
	} {
		assert.Equal(t, c.answer, testCompare(c.left, c.right))
	}
}

func TestStringToString(t *testing.T) {
	s := NewString("foo")
	assert.Equal(t, s, EvalPure(PApp(ToString, s)))
}
