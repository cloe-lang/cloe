package core

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSignatureBind(t *testing.T) {
	for _, c := range []struct {
		signature Signature
		arguments Arguments
	}{
		{
			NewSignature(nil, "", nil, ""),
			NewArguments(nil, nil),
		},
		{
			NewSignature([]string{"x"}, "", nil, ""),
			NewArguments([]PositionalArgument{NewPositionalArgument(NewList(Nil), true)}, nil),
		},
		{
			NewSignature(nil, "", []OptionalParameter{NewOptionalParameter("x", Nil)}, ""),
			NewArguments(nil, nil),
		},
		{
			NewSignature(nil, "", []OptionalParameter{NewOptionalParameter("foo", Nil)}, ""),
			NewArguments(nil, []KeywordArgument{NewKeywordArgument("", NewDictionary([]KeyValue{{NewString("foo"), Nil}}))}),
		},
		{
			NewSignature(nil, "", nil, "foo"),
			NewArguments(
				nil,
				[]KeywordArgument{NewKeywordArgument("foo", Nil), NewKeywordArgument("", NewDictionary([]KeyValue{{NewString("bar"), Nil}}))}),
		},
	} {
		vs, err := c.signature.Bind(c.arguments)
		assert.Equal(t, c.signature.arity(), len(vs))
		assert.Equal(t, nil, err)
	}
}

func TestSignatureBindError(t *testing.T) {
	for _, c := range []struct {
		signature Signature
		arguments Arguments
	}{
		{
			NewSignature([]string{"x"}, "", nil, ""),
			NewArguments(nil, nil),
		},
		{
			NewSignature(nil, "", nil, ""),
			NewArguments([]PositionalArgument{NewPositionalArgument(Nil, false)}, nil),
		},
		{
			NewSignature(nil, "", []OptionalParameter{NewOptionalParameter("arg", Nil)}, ""),
			NewArguments([]PositionalArgument{NewPositionalArgument(Nil, false)}, nil),
		},
	} {
		_, err := c.signature.Bind(c.arguments)
		assert.NotEqual(t, nil, err)
	}
}
