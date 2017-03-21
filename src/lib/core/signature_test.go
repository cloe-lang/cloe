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

func TestSignatureBindExpandedDictionaries(t *testing.T) {
	insert := NewLazyFunction(
		NewSignature(
			[]string{"collection", "key", "value"}, []OptionalArgument{}, "",
			[]string{}, []OptionalArgument{}, "",
		),
		func(ts ...*Thunk) (result Value) {
			return PApp(Insert, ts...)
		})

	f := App(Partial, NewArguments(
		[]PositionalArgument{
			NewPositionalArgument(insert, false),
			NewPositionalArgument(EmptyDictionary, false),
		},
		[]KeywordArgument{},
		[]*Thunk{NewDictionary([]Value{NewString("key").Eval()}, []*Thunk{True})}))

	v := App(f, NewArguments(
		[]PositionalArgument{},
		[]KeywordArgument{NewKeywordArgument("value", NewNumber(42))},
		[]*Thunk{})).Eval()

	_, ok := v.(DictionaryType)
	assert.True(t, ok)

	// Check if the Arguments passed to Partial is persistent.

	v = App(f, NewArguments(
		[]PositionalArgument{},
		[]KeywordArgument{NewKeywordArgument("value", NewNumber(42))},
		[]*Thunk{})).Eval()

	_, ok = v.(DictionaryType)
	assert.True(t, ok)
}
