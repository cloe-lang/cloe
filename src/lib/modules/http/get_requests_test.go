package http

import (
	"io/ioutil"
	"net/http"
	"testing"
	"time"

	"github.com/coel-lang/coel/src/lib/core"
	"github.com/coel-lang/coel/src/lib/systemt"
	"github.com/stretchr/testify/assert"
)

func TestGetRequests(t *testing.T) {
	go systemt.RunDaemons()

	th := core.PApp(getRequests, core.NewString(":8080"))

	go th.Eval()
	time.Sleep(100 * time.Millisecond)

	rc := make(chan string)
	go func() {
		r, err := http.Get("http://127.0.0.1:8080/foo/bar?baz=123")

		assert.Nil(t, err)

		bs, err := ioutil.ReadAll(r.Body)

		assert.Nil(t, err)

		rc <- string(bs)
	}()

	_, ok := th.Eval().(core.ListType)

	assert.True(t, ok)

	r := core.PApp(core.First, th)

	testRequest(t, r)

	th = core.PApp(
		core.PApp(r, core.NewString("respond")),
		core.NewString("Hello, world!"))

	assert.Equal(t, core.Nil.Eval(), core.PApp(core.Pure, th).Eval())

	assert.Equal(t, "Hello, world!", <-rc)
}

func testRequest(t *testing.T, th core.Value) {
	assert.Equal(t, core.NewString("").Eval(), core.PApp(th, core.NewString("body")).Eval())
	assert.Equal(t, core.NewString("GET").Eval(), core.PApp(th, core.NewString("method")).Eval())
	assert.Equal(t,
		core.NewString("/foo/bar?baz=123").Eval(),
		core.PApp(th, core.NewString("url")).Eval())
}

func TestGetRequestsWithCustomStatus(t *testing.T) {
	go systemt.RunDaemons()

	th := core.PApp(getRequests, core.NewString(":8888"))

	go th.Eval()
	time.Sleep(100 * time.Millisecond)

	status := make(chan int)
	go func() {
		r, err := http.Get("http://127.0.0.1:8888/foo/bar?baz=123")

		assert.Nil(t, err)

		status <- r.StatusCode
	}()

	th = core.App(
		core.PApp(core.PApp(core.First, th), core.NewString("respond")),
		core.NewArguments(
			nil,
			[]core.KeywordArgument{
				core.NewKeywordArgument("status", core.NewNumber(404)),
			},
			nil))

	assert.Equal(t, core.Nil.Eval(), core.PApp(core.Pure, th).Eval())
	assert.Equal(t, 404, <-status)
}
