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
	_, err := s.String("a")()
	assert.Equal(t, s.Line(), "a")
	assert.Equal(t, nil, err)
}

func TestStateLineAfterNewLine(t *testing.T) {
	s := NewState("a\nb\n")
	_, err := s.String("a\n")()
	assert.Equal(t, s.Line(), "b")
	assert.Equal(t, nil, err)
}

func TestStateLineNumber(t *testing.T) {
	assert.Equal(t, 1, NewState("").LineNumber())
}

func TestStateLinePosition(t *testing.T) {
	assert.Equal(t, 1, NewState("").LinePosition())
}
