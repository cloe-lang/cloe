package comb

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStateLine(t *testing.T) {
	assert.Equal(t, NewState("a\nb\n").Line(), "a")
}

func TestStateLineBeforeNewLine(t *testing.T) {
	s := NewState("a\nb\n")
	s.String("a")()
	assert.Equal(t, s.Line(), "a")
}

func TestStateLineAfterNewLine(t *testing.T) {
	s := NewState("a\nb\n")
	s.String("a\n")()
	assert.Equal(t, s.Line(), "b")
}
