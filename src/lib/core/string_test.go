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
	s := "foo"
	th := NewString(s)

	assert.Equal(t, string(PApp(Merge, th, th).Eval().(StringType)), s+s)
}

func TestStringToList(t *testing.T) {
	s := "lisp"
	l := PApp(ToList, NewString(s))

	for _, r := range s {
		assert.Equal(t, string(PApp(First, l).Eval().(StringType)), string(r))
		l = PApp(Rest, l)
	}

	assert.Equal(t, l.Eval().(ListType), emptyList)
}

func TestStringIndex(t *testing.T) {
	assert.Equal(t,
		NewString("b").Eval(),
		PApp(NewString("abc"), NewNumber(1)).Eval().(StringType))
}

func TestStringDelete(t *testing.T) {
	for _, test := range []struct {
		string   StringType
		index    float64
		expected StringType
	}{
		{"a", 0, ""},
		{"abc", 1, "ac"},
		{"abc", 2, "ab"},
	} {
		assert.Equal(t,
			test.expected,
			PApp(Delete, Normal(test.string), NewNumber(test.index)).Eval().(StringType))
	}
}

func TestStringSize(t *testing.T) {
	for _, test := range []struct {
		string StringType
		size   NumberType
	}{
		{"", 0},
		{"a", 1},
		{"ab", 2},
		{"abc", 3},
	} {
		assert.Equal(t, test.size, PApp(Size, Normal(test.string)).Eval().(NumberType))
	}
}

func TestStringInclude(t *testing.T) {
	for _, test := range []struct {
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
		assert.Equal(t, test.answer, PApp(Include, Normal(test.string), Normal(test.substring)).Eval().(BoolType))
	}
}

func TestStringInsert(t *testing.T) {
	for _, test := range []struct {
		string   StringType
		index    NumberType
		elem     StringType
		expected StringType
	}{
		{"", 0, "", ""},
		{"", 0, "a", "a"},
		{"a", 0, "b", "ba"},
		{"a", 1, "b", "ab"},
		{"ab", 0, "c", "cab"},
		{"ab", 1, "c", "acb"},
		{"ab", 2, "c", "abc"},
	} {
		assert.True(t, testEqual(
			Normal(test.expected),
			PApp(Insert, Normal(test.string), Normal(test.index), Normal(test.elem))))
	}
}

func TestStringCompare(t *testing.T) {
	assert.True(t, testCompare(NewString("foo"), NewString("foo")) == 0)
	assert.True(t, testCompare(NewString("foo"), NewString("bar")) == 1)
	assert.True(t, testCompare(NewString("bar"), NewString("foo")) == -1)
}
