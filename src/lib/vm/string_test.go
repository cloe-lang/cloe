package vm

import (
	"github.com/stretchr/testify/assert"
	"testing"
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
