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

	assert.Equal(t, string(App(Merge, th, NewList(th)).Eval().(StringType)), s+s)
}

func TestStringToList(t *testing.T) {
	s := "lisp"
	l := App(ToList, NewString(s))

	for _, r := range s {
		assert.Equal(t, string(App(First, l).Eval().(StringType)), string(r))
		l = App(Rest, l)
	}

	assert.Equal(t, l.Eval().(ListType), emptyList)
}
