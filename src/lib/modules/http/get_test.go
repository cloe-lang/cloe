package http

import (
	"strings"
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
	assert.True(t, strings.Contains(e.Error(), "TypeError"))
}

func TestGetWithInvalidURL(t *testing.T) {
	e, ok := core.PApp(get, core.NewString("http://hey-hey-i-am-invalid")).Eval().(core.ErrorType)

	assert.True(t, ok)
	assert.True(t, strings.Contains(e.Error(), "HTTPError"))
}
