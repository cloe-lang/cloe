package re

import (
	"testing"

	"github.com/coel-lang/coel/src/lib/core"
	"github.com/stretchr/testify/assert"
)

func TestFind(t *testing.T) {
	for _, c := range []struct {
		pattern, src string
		answer       *core.Thunk
	}{
		{"foo", "f", core.Nil},
		{"foo", "foo", core.NewList(core.NewString("foo"))},
		{"fo", "foo", core.NewList(core.NewString("fo"))},
		{"f(o)*", "f", core.NewList(core.NewString("f"), core.Nil)},
		{"f(o)*", "afoo", core.NewList(core.NewString("foo"), core.NewString("o"))},
	} {
		th := core.PApp(
			find,
			core.NewString(c.pattern),
			core.NewString(c.src))

		t.Log(core.PApp(core.Dump, th).Eval())

		b, ok := core.PApp(core.Equal, th, c.answer).Eval().(core.BoolType)

		assert.True(t, ok)
		assert.True(t, bool(b))
	}
}

func TestFindError(t *testing.T) {
	for _, th := range []*core.Thunk{
		core.PApp(find),
		core.PApp(find, core.NewString("foo")),
		core.PApp(find, core.NewString("foo"), core.Nil),
		core.PApp(find, core.NewString("(foo"), core.NewString("foo")),
	} {
		_, ok := th.Eval().(core.ErrorType)
		assert.True(t, ok)
	}
}
