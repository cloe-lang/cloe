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
			NewArguments(nil, nil, nil),
		},
		{
			NewSignature([]string{"x"}, "", nil, ""),
			NewArguments([]PositionalArgument{NewPositionalArgument(NewList(Nil), true)}, nil, nil),
		},
		{
			NewSignature(nil, "", []OptionalParameter{NewOptionalParameter("x", Nil)}, ""),
			NewArguments(nil, nil, nil),
		},
		{
			NewSignature(nil, "", []OptionalParameter{NewOptionalParameter("foo", Nil)}, ""),
			NewArguments(nil, nil, []Value{NewDictionary([]KeyValue{{NewString("foo"), Nil}})}),
		},
		{
			NewSignature(nil, "", nil, "foo"),
			NewArguments(
				nil,
				[]KeywordArgument{NewKeywordArgument("foo", Nil)},
				[]Value{NewDictionary([]KeyValue{{NewString("bar"), Nil}})}),
		},
	} {
		vs, err := c.signature.Bind(c.arguments)
		assert.Equal(t, c.signature.arity(), len(vs))
		assert.Equal(t, Value(nil), err)
	}
}

func TestSignatureBindError(t *testing.T) {
	for _, c := range []struct {
		signature Signature
		arguments Arguments
	}{
		{
			NewSignature(nil, "", nil, ""),
			NewArguments(nil, nil, []Value{Nil}),
		},
		{
			NewSignature([]string{"x"}, "", nil, ""),
			NewArguments(nil, nil, nil),
		},
		{
			NewSignature(nil, "", nil, ""),
			NewArguments([]PositionalArgument{NewPositionalArgument(Nil, false)}, nil, nil),
		},
		{
			NewSignature(nil, "", []OptionalParameter{NewOptionalParameter("arg", Nil)}, ""),
			NewArguments([]PositionalArgument{NewPositionalArgument(Nil, false)}, nil, nil),
		},
	} {
		_, err := c.signature.Bind(c.arguments)
		assert.NotEqual(t, Value(nil), err)
	}
}
