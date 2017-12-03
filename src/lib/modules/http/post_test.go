package http

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tisp-lang/tisp/src/lib/core"
)

func TestPost(t *testing.T) {
	th := core.PApp(post, core.NewString("https://google.com"), core.NewString(""))
	_, ok := th.Eval().(core.DictionaryType)

	assert.True(t, ok)
	assert.Equal(t,
		405.0,
		float64(core.PApp(core.Index, th, core.NewString("status")).Eval().(core.NumberType)))

	_, ok = core.PApp(core.Index, th, core.NewString("body")).Eval().(core.StringType)

	assert.True(t, ok)
}

func TestPostWithInvalidURL(t *testing.T) {
	e, ok := core.PApp(post, core.Nil, core.NewString("")).Eval().(core.ErrorType)

	assert.True(t, ok)
	assert.Equal(t, "TypeError", e.Name())
}

func TestPostWithInvalidBody(t *testing.T) {
	e, ok := core.PApp(post, core.NewString("https://google.com"), core.Nil).Eval().(core.ErrorType)

	assert.True(t, ok)
	assert.Equal(t, "TypeError", e.Name())
}
