package re

import (
	"testing"

	"github.com/coel-lang/coel/src/lib/core"
	"github.com/stretchr/testify/assert"
)

func TestReplace(t *testing.T) {
	for _, c := range []struct {
		pattern, repl, src, dest string
	}{
		{"foo", "bar", "foo", "bar"},
		{"f(.*)a", "b${1}r", "fooa", "boor"},
	} {
		s, ok := core.PApp(
			replace,
			core.NewString(c.pattern),
			core.NewString(c.repl),
			core.NewString(c.src)).Eval().(core.StringType)

		assert.True(t, ok)
		assert.Equal(t, c.dest, string(s))
	}
}

func TestReplaceError(t *testing.T) {
	for _, th := range []*core.Thunk{
		core.PApp(replace),
		core.PApp(replace, core.NewString("foo")),
		core.PApp(replace, core.NewString("foo"), core.Nil, core.NewString("bar")),
		core.PApp(replace, core.NewString("(foo"), core.NewString("foo"), core.NewString("foo")),
	} {
		_, ok := th.Eval().(core.ErrorType)
		assert.True(t, ok)
	}
}
