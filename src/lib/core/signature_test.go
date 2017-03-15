package core

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestXFailSignatureBindNoArgument(t *testing.T) {
	s := NewSignature([]string{}, []OptionalArgument{}, "", []string{}, []OptionalArgument{}, "")
	args := NewArguments([]PositionalArgument{NewPositionalArgument(Nil, false)}, []KeywordArgument{}, []*Thunk{})
	_, err := s.Bind(args)
	assert.NotEqual(t, (*Thunk)(nil), err)
}

func TestXFailSignatureBindRequiredKeywordArgument(t *testing.T) {
	s := NewSignature([]string{}, []OptionalArgument{}, "", []string{"arg"}, []OptionalArgument{}, "")
	args := NewArguments([]PositionalArgument{NewPositionalArgument(Nil, false)}, []KeywordArgument{}, []*Thunk{})
	_, err := s.Bind(args)
	assert.NotEqual(t, (*Thunk)(nil), err)
}

func TestXFailSignatureBindOptionalKeywordArgument(t *testing.T) {
	s := NewSignature([]string{}, []OptionalArgument{}, "", []string{}, []OptionalArgument{NewOptionalArgument("arg", Nil)}, "")
	args := NewArguments([]PositionalArgument{NewPositionalArgument(Nil, false)}, []KeywordArgument{}, []*Thunk{})
	_, err := s.Bind(args)
	assert.NotEqual(t, (*Thunk)(nil), err)
}
