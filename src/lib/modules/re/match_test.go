package re

import (
	"testing"

	"github.com/coel-lang/coel/src/lib/core"
	"github.com/stretchr/testify/assert"
)

func TestMatch(t *testing.T) {
	for _, c := range []struct {
		pattern, src string
		answer       bool
	}{
		{"foo", "foo", true},
		{"fo", "foo", true},
		{"foo", "fo", false},
		{"f(o)*", "f", true},
		{"f(o)*", "afoo", true},
	} {
		b, ok := core.PApp(
			match,
			core.NewString(c.pattern),
			core.NewString(c.src)).Eval().(core.BoolType)

		assert.True(t, ok)
		assert.Equal(t, c.answer, bool(b))
	}
}

func TestMatchError(t *testing.T) {
	for _, th := range []core.Value{
		core.PApp(match),
		core.PApp(match, core.NewString("foo")),
		core.PApp(match, core.NewString("foo"), core.Nil),
		core.PApp(match, core.NewString("(foo"), core.NewString("foo")),
	} {
		_, ok := th.Eval().(core.ErrorType)
		assert.True(t, ok)
	}
}
