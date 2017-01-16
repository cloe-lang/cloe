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

func TestStringAdd(t *testing.T) {
	s := "foo"
	th := NewString(s)

	assert.Equal(t, string(App(Add, th, th).Eval().(stringType)), s+s)
}
