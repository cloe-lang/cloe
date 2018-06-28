package http

import (
	"net/http"
	"os"
	"testing"
	"time"
)

const (
	rootURL        = "http://localhost:2049/"
	postURL        = "http://localhost:2049/post"
	invalidPathURL = "http://localhost:2049/invalid-path"
	invalidHostURL = "http://hostlocal:2049/"
)

type testHandler struct{}

func (testHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "/":
		w.Write(nil)
	case "/post":
		if r.Method != "POST" {
			w.WriteHeader(404)
		}
	default:
		w.WriteHeader(404)
	}
}

func TestMain(m *testing.M) {
	go http.ListenAndServe(":2049", testHandler{})

	time.Sleep(time.Millisecond)

	os.Exit(m.Run())
}
