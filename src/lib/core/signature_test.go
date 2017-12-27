package core

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSignatureBindWithEmtpyArguments(t *testing.T) {
	s := NewSignature([]string{"x"}, nil, "", nil, nil, "")
	args := NewArguments(nil, nil, nil)
	_, err := s.Bind(args)
	assert.NotEqual(t, (*Thunk)(nil), err)
}

func TestSignatureBindWithOptionalPositionalArguments(t *testing.T) {
	s := NewSignature(nil, []OptionalArgument{NewOptionalArgument("x", Nil)}, "", nil, nil, "")

	for _, args := range []Arguments{
		NewArguments(nil, nil, nil),
		NewArguments([]PositionalArgument{NewPositionalArgument(True, false)}, nil, nil),
	} {
		_, err := s.Bind(args)
		assert.Equal(t, (*Thunk)(nil), err)
	}
}

func TestSignatureBindExpandedPositionalArgument(t *testing.T) {
	s := NewSignature([]string{"x"}, nil, "", nil, nil, "")
	args := NewArguments([]PositionalArgument{NewPositionalArgument(NewList(Nil), true)}, nil, nil)
	_, err := s.Bind(args)
	assert.Equal(t, (*Thunk)(nil), err)
}

func TestSignatureBindExpandedDictionary(t *testing.T) {
	s := NewSignature(nil, nil, "", []string{"foo"}, nil, "")
	args := NewArguments(
		nil,
		nil,
		[]*Thunk{NewDictionary([]KeyValue{{NewString("foo"), Nil}})})
	_, err := s.Bind(args)
	assert.Equal(t, (*Thunk)(nil), err)
}

func TestSignatureBindExpandedDictionaryFail(t *testing.T) {
	s := NewSignature(nil, nil, "", []string{"foo"}, nil, "")
	args := NewArguments(nil, nil, []*Thunk{Nil})
	_, err := s.Bind(args)
	assert.NotEqual(t, (*Thunk)(nil), err)
}

func TestSignatureBindRestKeywordArguments(t *testing.T) {
	s := NewSignature(nil, nil, "", nil, nil, "foo")
	args := NewArguments(
		nil,
		[]KeywordArgument{NewKeywordArgument("foo", Nil)},
		[]*Thunk{NewDictionary([]KeyValue{{NewString("bar"), Nil}})})
	_, err := s.Bind(args)
	assert.Equal(t, (*Thunk)(nil), err)
}

func TestSignatureBindNoArgumentFail(t *testing.T) {
	s := NewSignature(nil, nil, "", nil, nil, "")
	args := NewArguments([]PositionalArgument{NewPositionalArgument(Nil, false)}, nil, nil)
	_, err := s.Bind(args)
	assert.NotEqual(t, (*Thunk)(nil), err)
}

func TestSignatureBindRequiredKeywordArgumentFail(t *testing.T) {
	s := NewSignature(nil, nil, "", []string{"arg"}, nil, "")
	args := NewArguments([]PositionalArgument{NewPositionalArgument(Nil, false)}, nil, nil)
	_, err := s.Bind(args)
	assert.NotEqual(t, (*Thunk)(nil), err)
}

func TestSignatureBindOptionalKeywordArgumentFail(t *testing.T) {
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
		[]*Thunk{NewDictionary([]KeyValue{{NewString("key"), True}})}))

	v := App(f, NewArguments(nil, []KeywordArgument{NewKeywordArgument("value", NewNumber(42))}, nil)).Eval()

	_, ok := v.(DictionaryType)
	assert.True(t, ok)

	// Check if the Arguments passed to Partial is persistent.

	v = App(f, NewArguments(nil, []KeywordArgument{NewKeywordArgument("value", NewNumber(42))}, nil)).Eval()

	_, ok = v.(DictionaryType)
	assert.True(t, ok)
}
