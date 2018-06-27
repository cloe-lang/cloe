package json

import (
	"testing"

	"github.com/cloe-lang/cloe/src/lib/core"
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
			core.NewDictionary([]core.KeyValue{{Key: core.Nil, Value: core.Nil}}),
			`{"null":null}`,
		},
		{
			core.NewDictionary([]core.KeyValue{{Key: core.NewNumber(42), Value: core.NewNumber(42)}}),
			`{"42":42}`,
		},
		{
			core.NewDictionary([]core.KeyValue{{Key: core.True, Value: core.True}}),
			`{"true":true}`,
		},
	} {
		s, ok := core.EvalPure(core.PApp(encode, c.value)).(core.StringType)

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
			{Key: core.NewString("foo"), Value: core.NewNumber(42)},
			{Key: core.NewString("bar"), Value: core.True},
			{Key: core.NewString("baz"), Value: core.NewString("blah")},
		}),
	} {
		s, ok := core.EvalPure(core.PApp(encode, th)).(core.StringType)

		assert.True(t, ok)

		b, ok := core.EvalPure(core.PApp(core.Equal, th, core.PApp(decode, s))).(*core.BooleanType)

		assert.True(t, ok)
		assert.True(t, bool(*b))
	}
}

func TestEncodeWithInvalidArguments(t *testing.T) {
	for _, th := range []core.Value{
		core.If,
		core.DummyError,
		core.PApp(core.Prepend, core.Nil, core.DummyError),
		core.NewList(core.DummyError),
		core.PApp(core.Insert, core.EmptyDictionary, core.NewString("foo"), core.DummyError),
		core.NewDictionary([]core.KeyValue{{Key: core.NewList(core.DummyError), Value: core.Nil}}),
		core.NewDictionary([]core.KeyValue{{Key: core.Nil, Value: core.DummyError}}),
	} {
		_, ok := core.EvalPure(core.PApp(encode, th)).(*core.ErrorType)
		assert.True(t, ok)
	}
}
