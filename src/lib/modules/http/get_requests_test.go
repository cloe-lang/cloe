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
		r, err := http.Get("http://127.0.0.1:8080")

		assert.Equal(t, nil, err)

		bs, err := ioutil.ReadAll(r.Body)

		assert.Equal(t, nil, err)

		rc <- string(bs)
	}()

	_, ok := th.Eval().(core.ListType)

	assert.True(t, ok)

	th = core.PApp(
		core.PApp(core.PApp(th, core.NewNumber(0)), core.NewString("sendResponse")),
		core.NewString("Hello, world!"))

	assert.Equal(t, core.Nil.Eval(), core.PApp(core.Pure, th).Eval())

	assert.Equal(t, "Hello, world!", <-rc)
}
