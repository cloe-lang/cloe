package http

import (
	"testing"

	"github.com/cloe-lang/cloe/src/lib/core"
	"github.com/stretchr/testify/assert"
)

func TestGet(t *testing.T) {
	v := core.PApp(get, core.NewString(rootURL))
	_, ok := core.EvalPure(v).(*core.DictionaryType)

	assert.True(t, ok)
	assert.Equal(t,
		core.NewNumber(200),
		core.EvalPure(core.PApp(core.Index, v, core.NewString("status"))))

	_, ok = core.EvalPure(core.PApp(core.Index, v, core.NewString("body"))).(core.StringType)

	assert.True(t, ok)
}

func TestGetWithInvalidArgument(t *testing.T) {
	e, ok := core.EvalPure(core.PApp(get, core.Nil)).(*core.ErrorType)

	assert.True(t, ok)
	assert.Equal(t, "TypeError", e.Name())
}

func TestGetWithInvalidHost(t *testing.T) {
	e, ok := core.EvalPure(core.PApp(get, core.NewString(invalidHostURL))).(*core.ErrorType)

	assert.True(t, ok)
	assert.Equal(t, "HTTPError", e.Name())
}

func TestGetWithInvalidPath(t *testing.T) {
	e, ok := core.EvalPure(core.PApp(get, core.NewString(invalidPathURL))).(*core.ErrorType)

	assert.True(t, ok)
	assert.Equal(t, "HTTPError", e.Name())
}

func TestGetWithInvalidPathButNoError(t *testing.T) {
	v := core.App(
		get,
		core.NewArguments(
			[]core.PositionalArgument{
				core.NewPositionalArgument(core.NewString(invalidPathURL), false),
			},
			[]core.KeywordArgument{
				core.NewKeywordArgument("error", core.False)}))

	_, ok := core.EvalPure(v).(*core.DictionaryType)
	assert.True(t, ok)

	n, err := core.EvalNumber(core.PApp(core.Index, v, core.NewString("status")))

	assert.Equal(t, nil, err)
	assert.Equal(t, 404.0, float64(n))
}

func TestGetWithInvalidErrorArgument(t *testing.T) {
	e, ok := core.EvalPure(core.App(
		get,
		core.NewArguments(
			[]core.PositionalArgument{
				core.NewPositionalArgument(core.NewString(invalidPathURL), false),
			},
			[]core.KeywordArgument{
				core.NewKeywordArgument("error", core.Nil)}))).(*core.ErrorType)

	assert.True(t, ok)
	assert.Equal(t, "TypeError", e.Name())
}
