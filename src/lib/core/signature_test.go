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
			NewArguments(nil, nil, []*Thunk{NewDictionary([]KeyValue{{NewString("foo"), Nil}})}),
		},
		{
			NewSignature(nil, nil, "", nil, nil, "foo"),
			NewArguments(
				nil,
				[]KeywordArgument{NewKeywordArgument("foo", Nil)},
				[]*Thunk{NewDictionary([]KeyValue{{NewString("bar"), Nil}})})},
	} {
		ts, err := c.signature.Bind(c.arguments)
		assert.Equal(t, c.signature.arity(), len(ts))
		assert.Equal(t, (*Thunk)(nil), err)
	}
}

func TestSignatureBindError(t *testing.T) {
	for _, c := range []struct {
		signature Signature
		arguments Arguments
	}{
		{
			NewSignature(nil, nil, "", []string{"foo"}, nil, ""),
			NewArguments(nil, nil, []*Thunk{Nil}),
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
		assert.NotEqual(t, (*Thunk)(nil), err)
	}
}

func TestSignatureBindExpandedDictionaries(t *testing.T) {
	insert := NewLazyFunction(
		NewSignature(
			[]string{"collection", "key", "value"}, nil, "",
			nil, nil, "",
		),
		func(ts ...*Thunk) (result Value) {
			return PApp(Insert, ts...)
		})

	f := App(Partial, NewArguments(
		[]PositionalArgument{
			NewPositionalArgument(insert, false),
			NewPositionalArgument(EmptyDictionary, false),
		},
		nil,
		[]*Thunk{NewDictionary([]KeyValue{{NewString("key"), True}})}))

	v := App(f, NewArguments(nil, []KeywordArgument{NewKeywordArgument("value", NewNumber(42))}, nil)).Eval()

	_, ok := v.(DictionaryType)
	assert.True(t, ok)

	// Check if the Arguments passed to Partial is persistent.

	v = App(f, NewArguments(nil, []KeywordArgument{NewKeywordArgument("value", NewNumber(42))}, nil)).Eval()

	_, ok = v.(DictionaryType)
	assert.True(t, ok)
}
