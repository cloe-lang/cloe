package http

import (
	"io/ioutil"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/tisp-lang/tisp/src/lib/core"
	"github.com/tisp-lang/tisp/src/lib/systemt"
)

func TestGetRequests(t *testing.T) {
	go systemt.RunDaemons()

	th := core.PApp(getRequests, core.NewString(":8080"))

	go th.Eval()
	time.Sleep(100 * time.Millisecond)

	rc := make(chan string)
	go func() {
		r, err := http.Get("http://127.0.0.1:8080/foo/bar?baz=123")

		assert.Equal(t, nil, err)

		bs, err := ioutil.ReadAll(r.Body)

		assert.Equal(t, nil, err)

		rc <- string(bs)
	}()

	_, ok := th.Eval().(core.ListType)

	assert.True(t, ok)

	r := core.PApp(th, core.NewNumber(0))

	testRequest(t, r)

	th = core.PApp(
		core.PApp(r, core.NewString("respond")),
		core.NewString("Hello, world!"))

	assert.Equal(t, core.Nil.Eval(), core.PApp(core.Pure, th).Eval())

	assert.Equal(t, "Hello, world!", <-rc)
}

func testRequest(t *testing.T, th *core.Thunk) {
	assert.Equal(t, core.NewString("").Eval(), core.PApp(th, core.NewString("body")).Eval())
	assert.Equal(t, core.NewString("GET").Eval(), core.PApp(th, core.NewString("method")).Eval())
	assert.Equal(t,
		core.NewString("/foo/bar?baz=123").Eval(),
		core.PApp(th, core.NewString("url")).Eval())
}
