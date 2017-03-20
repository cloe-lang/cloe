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
