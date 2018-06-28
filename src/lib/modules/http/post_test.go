package http

import (
	"testing"

	"github.com/cloe-lang/cloe/src/lib/core"
	"github.com/stretchr/testify/assert"
)

func TestPost(t *testing.T) {
	v := core.PApp(post, core.NewString(postURL), core.NewString(""))
	_, ok := core.EvalPure(v).(*core.DictionaryType)

	t.Log(core.EvalPure(v))
	assert.True(t, ok)
	assert.Equal(t,
		core.NewNumber(200),
		core.EvalPure(core.PApp(core.Index, v, core.NewString("status"))))

	_, ok = core.EvalPure(core.PApp(core.Index, v, core.NewString("body"))).(core.StringType)

	assert.True(t, ok)
}

func TestPostWithInvalidURL(t *testing.T) {
	e, ok := core.EvalPure(core.PApp(post, core.Nil, core.NewString(""))).(*core.ErrorType)

	assert.True(t, ok)
	assert.Equal(t, "TypeError", e.Name())
}

func TestPostWithInvalidBody(t *testing.T) {
	e, ok := core.EvalPure(core.PApp(post, core.NewString(rootURL), core.Nil)).(*core.ErrorType)

	assert.True(t, ok)
	assert.Equal(t, "TypeError", e.Name())
}
