package json

import (
	"testing"

	"github.com/coel-lang/coel/src/lib/core"
	"github.com/stretchr/testify/assert"
)

func TestEncode(t *testing.T) {
	for _, c := range []struct {
		value  core.Value
		answer string
	}{
		{core.True, `true`},
		{core.False, `false`},
		{core.NewNumber(123), `123`},
		{core.NewNumber(-0.125), `-0.125`},
		{core.NewString("foo"), `"foo"`},
		{core.Nil, `null`},
		{core.EmptyList, `[]`},
		{core.NewList(core.If), `[]`},
		{core.NewList(
			core.NewString("foo"),
			core.NewNumber(42),
			core.NewString("bar"),
		), `["foo",42,"bar"]`},
		{core.EmptyDictionary, "{}"},
		{
			core.NewDictionary([]core.KeyValue{{core.Nil, core.Nil}}),
			`{"null":null}`,
		},
		{
			core.NewDictionary([]core.KeyValue{{core.NewNumber(42), core.NewNumber(42)}}),
			`{"42":42}`,
		},
		{
			core.NewDictionary([]core.KeyValue{{core.True, core.True}}),
			`{"true":true}`,
		},
	} {
		s, ok := core.PApp(encode, c.value).Eval().(core.StringType)

		assert.True(t, ok)
		assert.Equal(t, c.answer, string(s))
	}
}

func TestEncodeAndDecode(t *testing.T) {
	for _, th := range []core.Value{
		core.True,
		core.False,
		core.NewNumber(123),
		core.NewNumber(-0.125),
		core.NewString("foo"),
		core.Nil,
		core.EmptyList,
		core.NewList(
			core.NewString("foo"),
			core.NewNumber(42),
			core.NewString("bar"),
		),
		core.EmptyDictionary,
		core.NewDictionary([]core.KeyValue{
			{core.NewString("foo"), core.NewNumber(42)},
			{core.NewString("bar"), core.True},
			{core.NewString("baz"), core.NewString("blah")},
		}),
	} {
		s, ok := core.PApp(encode, th).Eval().(core.StringType)

		assert.True(t, ok)

		b, ok := core.PApp(core.Equal, th, core.PApp(decode, s)).Eval().(core.BoolType)

		assert.True(t, ok)
		assert.True(t, bool(b))
	}
}

func TestEncodeWithInvalidArguments(t *testing.T) {
	for _, th := range []core.Value{
		core.If,
		core.ValueError("It's wrong."),
		core.PApp(core.Prepend, core.Nil, core.ValueError("Not a list")),
		core.NewList(core.ValueError("")),
		core.PApp(core.Insert, core.EmptyDictionary, core.NewString("foo"), core.ValueError("")),
		core.NewDictionary([]core.KeyValue{{core.NewList(core.DummyError), core.Nil}}),
		core.NewDictionary([]core.KeyValue{{core.Nil, core.DummyError}}),
	} {
		_, ok := core.PApp(encode, th).Eval().(core.ErrorType)
		assert.True(t, ok)
	}
}
