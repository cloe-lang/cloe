package re

import (
	"testing"

	"github.com/cloe-lang/cloe/src/lib/core"
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
		b, ok := core.EvalPure(core.PApp(
			match,
			core.NewString(c.pattern),
			core.NewString(c.src))).(*core.BooleanType)

		assert.True(t, ok)
		assert.Equal(t, c.answer, bool(*b))
	}
}

func TestMatchError(t *testing.T) {
	for _, v := range []core.Value{
		core.PApp(match),
		core.PApp(match, core.NewString("foo")),
		core.PApp(match, core.NewString("foo"), core.Nil),
		core.PApp(match, core.NewString("(foo"), core.NewString("foo")),
	} {
		_, ok := core.EvalPure(v).(*core.ErrorType)
		assert.True(t, ok)
	}
}
