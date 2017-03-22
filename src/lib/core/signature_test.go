package core

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestXFailSignatureBindNoArgument(t *testing.T) {
	s := NewSignature(nil, nil, "", nil, nil, "")
	args := NewArguments([]PositionalArgument{NewPositionalArgument(Nil, false)}, nil, []*Thunk{})
	_, err := s.Bind(args)
	assert.NotEqual(t, (*Thunk)(nil), err)
}

func TestXFailSignatureBindRequiredKeywordArgument(t *testing.T) {
	s := NewSignature(nil, nil, "", []string{"arg"}, nil, "")
	args := NewArguments([]PositionalArgument{NewPositionalArgument(Nil, false)}, nil, []*Thunk{})
	_, err := s.Bind(args)
	assert.NotEqual(t, (*Thunk)(nil), err)
}

func TestXFailSignatureBindOptionalKeywordArgument(t *testing.T) {
	s := NewSignature(nil, nil, "", nil, []OptionalArgument{NewOptionalArgument("arg", Nil)}, "")
	args := NewArguments([]PositionalArgument{NewPositionalArgument(Nil, false)}, nil, nil)
	_, err := s.Bind(args)
	assert.NotEqual(t, (*Thunk)(nil), err)
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
		[]*Thunk{NewDictionary([]Value{NewString("key").Eval()}, []*Thunk{True})}))

	v := App(f, NewArguments(nil, []KeywordArgument{NewKeywordArgument("value", NewNumber(42))}, nil)).Eval()

	_, ok := v.(DictionaryType)
	assert.True(t, ok)

	// Check if the Arguments passed to Partial is persistent.

	v = App(f, NewArguments(nil, []KeywordArgument{NewKeywordArgument("value", NewNumber(42))}, nil)).Eval()

	_, ok = v.(DictionaryType)
	assert.True(t, ok)
}
