package http

import (
	"testing"

	"github.com/coel-lang/coel/src/lib/core"
	"github.com/stretchr/testify/assert"
)

func TestGet(t *testing.T) {
	v := core.PApp(get, core.NewString("http://httpbin.org"))
	_, ok := core.EvalPure(v).(*core.DictionaryType)

	assert.True(t, ok)
	assert.Equal(t, core.NewNumber(200), core.EvalPure(core.PApp(v, core.NewString("status"))))

	_, ok = core.EvalPure(core.PApp(v, core.NewString("body"))).(*core.StringType)

	assert.True(t, ok)
}

func TestGetWithInvalidArgument(t *testing.T) {
	e, ok := core.EvalPure(core.PApp(get, core.Nil)).(*core.ErrorType)

	assert.True(t, ok)
	assert.Equal(t, "TypeError", e.Name())
}

func TestGetWithInvalidHost(t *testing.T) {
	e, ok := core.EvalPure(core.PApp(get, core.NewString("http://hey-hey-i-am-invalid"))).(*core.ErrorType)

	assert.True(t, ok)
	assert.Equal(t, "HTTPError", e.Name())
}

func TestGetWithInvalidPath(t *testing.T) {
	e, ok := core.EvalPure(core.PApp(get, core.NewString("http://httpbin.org/invalid-path"))).(*core.ErrorType)

	assert.True(t, ok)
	assert.Equal(t, "HTTPError", e.Name())
}

func TestGetWithInvalidPathButNoError(t *testing.T) {
	v := core.App(
		get,
		core.NewArguments(
			[]core.PositionalArgument{
				core.NewPositionalArgument(core.NewString("http://httpbin.org/invalid-path"), false),
			},
			[]core.KeywordArgument{
				core.NewKeywordArgument("error", core.False)}))

	_, ok := core.EvalPure(v).(*core.DictionaryType)
	assert.True(t, ok)

	n, err := core.EvalNumber(core.PApp(v, core.NewString("status")))

	assert.Equal(t, nil, err)
	assert.Equal(t, 404.0, float64(n))
}

func TestGetWithInvalidErrorArgument(t *testing.T) {
	e, ok := core.EvalPure(core.App(
		get,
		core.NewArguments(
			[]core.PositionalArgument{
				core.NewPositionalArgument(core.NewString("http://httpbin.org/invalid-path"), false),
			},
			[]core.KeywordArgument{
				core.NewKeywordArgument("error", core.Nil)}))).(*core.ErrorType)

	assert.True(t, ok)
	assert.Equal(t, "TypeError", e.Name())
}
