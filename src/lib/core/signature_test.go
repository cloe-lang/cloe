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
			NewSignature(nil, nil, "", nil, nil, ""),
			NewArguments(nil, nil, nil),
		},
		{
			NewSignature([]string{"x"}, nil, "", nil, nil, ""),
			NewArguments([]PositionalArgument{NewPositionalArgument(NewList(Nil), true)}, nil, nil),
		},
		{
			NewSignature(nil, []OptionalArgument{NewOptionalArgument("x", Nil)}, "", nil, nil, ""),
			NewArguments(nil, nil, nil),
		},
		{
			NewSignature(nil, []OptionalArgument{NewOptionalArgument("x", Nil)}, "", nil, nil, ""),
			NewArguments([]PositionalArgument{NewPositionalArgument(True, false)}, nil, nil),
		},
		{
			NewSignature(nil, nil, "", []string{"foo"}, nil, ""),
			NewArguments(nil, nil, []Value{NewDictionary([]KeyValue{{NewString("foo"), Nil}})}),
		},
		{
			NewSignature(nil, nil, "", nil, nil, "foo"),
			NewArguments(
				nil,
				[]KeywordArgument{NewKeywordArgument("foo", Nil)},
				[]Value{NewDictionary([]KeyValue{{NewString("bar"), Nil}})})},
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
			NewSignature(nil, nil, "", []string{"foo"}, nil, ""),
			NewArguments(nil, nil, []Value{Nil}),
		},
		{
			NewSignature([]string{"x"}, nil, "", nil, nil, ""),
			NewArguments(nil, nil, nil),
		},
		{
			NewSignature(nil, nil, "", nil, nil, ""),
			NewArguments([]PositionalArgument{NewPositionalArgument(Nil, false)}, nil, nil),
		},
		{
			NewSignature(nil, nil, "", []string{"arg"}, nil, ""),
			NewArguments([]PositionalArgument{NewPositionalArgument(Nil, false)}, nil, nil),
		},
		{
			NewSignature(nil, nil, "", nil, []OptionalArgument{NewOptionalArgument("arg", Nil)}, ""),
			NewArguments([]PositionalArgument{NewPositionalArgument(Nil, false)}, nil, nil),
		},
	} {
		_, err := c.signature.Bind(c.arguments)
		assert.NotEqual(t, Value(nil), err)
	}
}

func TestSignatureBindExpandedDictionaries(t *testing.T) {
	insert := NewLazyFunction(
		NewSignature(
			[]string{"collection", "key", "value"}, nil, "",
			nil, nil, "",
		),
		func(vs ...Value) (result Value) {
			return PApp(Insert, vs...)
		})

	f := App(Partial, NewArguments(
		[]PositionalArgument{
			NewPositionalArgument(insert, false),
			NewPositionalArgument(EmptyDictionary, false),
		},
		nil,
		[]Value{NewDictionary([]KeyValue{{NewString("key"), True}})}))

	v := EvalPure(App(f, NewArguments(nil, []KeywordArgument{NewKeywordArgument("value", NewNumber(42))}, nil)))

	_, ok := v.(DictionaryType)
	assert.True(t, ok)

	// Check if the Arguments passed to Partial is persistent.

	v = EvalPure(App(f, NewArguments(nil, []KeywordArgument{NewKeywordArgument("value", NewNumber(42))}, nil)))

	_, ok = v.(DictionaryType)
	assert.True(t, ok)
}
