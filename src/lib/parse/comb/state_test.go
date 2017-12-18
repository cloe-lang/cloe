package comb

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStateLine(t *testing.T) {
	assert.Equal(t, "a", NewState("a\nb\n").Line())
}

func TestStateLineBeforeNewLine(t *testing.T) {
	s := NewState("a\nb\n")
	_, err := s.String("a")()
	assert.Equal(t, "a", s.Line())
	assert.Nil(t, err)
}

func TestStateLineAfterNewLine(t *testing.T) {
	s := NewState("a\nb\n")
	_, err := s.String("a\n")()
	assert.Equal(t, "b", s.Line())
	assert.Nil(t, err)
}

func TestStateLineNumber(t *testing.T) {
	assert.Equal(t, 1, NewState("").LineNumber())
}

func TestStateLinePosition(t *testing.T) {
	assert.Equal(t, 1, NewState("").LinePosition())
}
