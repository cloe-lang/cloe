package http

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tisp-lang/tisp/src/lib/core"
)

func TestGet(t *testing.T) {
	th := core.PApp(get, core.NewString("https://google.com"))
	_, ok := th.Eval().(core.DictionaryType)

	assert.True(t, ok)
	assert.Equal(t, 200.0, float64(core.PApp(core.Index, th, core.NewString("status")).Eval().(core.NumberType)))

	_, ok = core.PApp(core.Index, th, core.NewString("body")).Eval().(core.StringType)

	assert.True(t, ok)
}

func TestGetWithInvalidArgument(t *testing.T) {
	e, ok := core.PApp(get, core.Nil).Eval().(core.ErrorType)

	assert.True(t, ok)
	assert.Equal(t, "TypeError", e.Name())
}

func TestGetWithInvalidHost(t *testing.T) {
	e, ok := core.PApp(get, core.NewString("http://hey-hey-i-am-invalid")).Eval().(core.ErrorType)

	assert.True(t, ok)
	assert.Equal(t, "HTTPError", e.Name())
}

func TestGetWithInvalidPath(t *testing.T) {
	e, ok := core.PApp(get, core.NewString("https://google.com/hey-google")).Eval().(core.ErrorType)

	assert.True(t, ok)
	assert.Equal(t, "HTTPError", e.Name())
}

func TestGetWithInvalidPathButNoError(t *testing.T) {
	th := core.App(
		get,
		core.NewArguments(
			[]core.PositionalArgument{
				core.NewPositionalArgument(core.NewString("https://google.com/hey-google"), false),
			},
			[]core.KeywordArgument{
				core.NewKeywordArgument("error", core.False)},
			nil))

	_, ok := th.Eval().(core.DictionaryType)
	assert.True(t, ok)

	n, ok := core.PApp(core.Index, th, core.NewString("status")).Eval().(core.NumberType)

	assert.True(t, ok)
	assert.Equal(t, 404.0, float64(n))
}

func TestGetWithInvalidErrorArgument(t *testing.T) {
	e, ok := core.App(
		get,
		core.NewArguments(
			[]core.PositionalArgument{
				core.NewPositionalArgument(core.NewString("https://google.com/hey-google"), false),
			},
			[]core.KeywordArgument{
				core.NewKeywordArgument("error", core.Nil)},
			nil)).Eval().(core.ErrorType)

	assert.True(t, ok)
	assert.Equal(t, "TypeError", e.Name())
}
