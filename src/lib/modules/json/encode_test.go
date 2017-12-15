package json

import (
	"testing"

	"github.com/coel-lang/coel/src/lib/core"
	"github.com/stretchr/testify/assert"
)

func TestEncode(t *testing.T) {
	for _, c := range []struct {
		value  *core.Thunk
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
			core.NewDictionary([]core.Value{core.Nil}, []*core.Thunk{core.Nil}),
			`{"null":null}`,
		},
		{
			core.NewDictionary([]core.Value{core.NewNumber(42)}, []*core.Thunk{core.NewNumber(42)}),
			`{"42":42}`,
		},
		{
			core.NewDictionary([]core.Value{core.True}, []*core.Thunk{core.True}),
			`{"true":true}`,
		},
	} {
		s, ok := core.PApp(encode, c.value).Eval().(core.StringType)

		assert.True(t, ok)
		assert.Equal(t, c.answer, string(s))
	}
}

func TestEncodeAndDecode(t *testing.T) {
	for _, th := range []*core.Thunk{
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
		core.NewDictionary(
			[]core.Value{
				core.NewString("foo").Eval(),
				core.NewString("bar").Eval(),
				core.NewString("baz").Eval(),
			},
			[]*core.Thunk{
				core.NewNumber(42),
				core.True,
				core.NewString("blah"),
			},
		),
	} {
		s, ok := core.PApp(encode, th).Eval().(core.StringType)

		assert.True(t, ok)

		b, ok := core.PApp(core.Equal, th, core.PApp(decode, core.Normal(s))).Eval().(core.BoolType)

		assert.True(t, ok)
		assert.True(t, bool(b))
	}
}

func TestEncodeWithInvalidArguments(t *testing.T) {
	for _, th := range []*core.Thunk{
		core.If,
		core.ValueError("It's wrong."),
		core.PApp(core.Prepend, core.Nil, core.ValueError("Not a list")),
		core.NewList(core.ValueError("")),
		core.PApp(core.Insert, core.EmptyDictionary, core.NewString("foo"), core.ValueError("")),
	} {
		_, ok := core.PApp(encode, th).Eval().(core.ErrorType)
		assert.True(t, ok)
	}
}
