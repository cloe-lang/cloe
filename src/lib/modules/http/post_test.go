package http

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/coel-lang/coel/src/lib/core"
)

func TestPost(t *testing.T) {
	th := core.PApp(post, core.NewString("http://httpbin.org/post"), core.NewString(""))
	_, ok := th.Eval().(core.DictionaryType)

	t.Log(th.Eval())
	assert.True(t, ok)
	assert.Equal(t,
		200.0,
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
	e, ok := core.PApp(post, core.NewString("http://httpbin.org"), core.Nil).Eval().(core.ErrorType)

	assert.True(t, ok)
	assert.Equal(t, "TypeError", e.Name())
}
