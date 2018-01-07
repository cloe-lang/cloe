package http

import (
	"io/ioutil"
	"net/http"
	"sync"

	"github.com/coel-lang/coel/src/lib/builtins"
	"github.com/coel-lang/coel/src/lib/core"
	"github.com/coel-lang/coel/src/lib/systemt"
)

const requestChannelSize = 1024
const responseChannelSize = 1024

var getRequests = core.NewLazyFunction(
	core.NewSignature([]string{"address"}, nil, "", nil, nil, ""),
	func(ts ...*core.Thunk) core.Value {
		s, err := ts[0].EvalString()

		if err != nil {
			return err
		}

		ec := make(chan error)
		h := newHandler()

		systemt.Daemonize(func() {
			if err := http.ListenAndServe(string(s), h); err != nil {
				ec <- err
			}
		})

		return core.PApp(core.PApp(builtins.Y, core.NewLazyFunction(
			core.NewSignature([]string{"me"}, nil, "", nil, nil, ""),
			func(ts ...*core.Thunk) core.Value {
				select {
				case t := <-h.Requests:
					return core.StrictPrepend([]*core.Thunk{t}, core.PApp(ts[0]))
				case err := <-ec:
					return httpError(err)
				}
			})))
	})

type handler struct {
	Requests  chan *core.Thunk
	responses <-chan string
}

func newHandler() handler {
	return handler{
		make(chan *core.Thunk, requestChannelSize),
		make(chan string, responseChannelSize),
	}
}

func (h handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	b, err := ioutil.ReadAll(r.Body)

	if err != nil {
		h.Requests <- httpError(err)
		return
	}

	wg := sync.WaitGroup{}

	wg.Add(1)

	h.Requests <- core.NewDictionary(
		[]core.KeyValue{
			{core.NewString("body"), core.NewString(string(b))},
			{core.NewString("method"), core.NewString(r.Method)},
			{core.NewString("url"), core.NewString(r.URL.String())},
			{
				core.NewString("respond"),
				core.NewEffectFunction(
					core.NewSignature(
						nil,
						[]core.OptionalArgument{
							core.NewOptionalArgument("body", core.NewString("")),
						}, "",
						nil,
						[]core.OptionalArgument{
							core.NewOptionalArgument("status", core.NewNumber(200)),
						},
						"",
					),
					func(ts ...*core.Thunk) core.Value {
						defer wg.Done()

						n, err := ts[1].EvalNumber()

						if err != nil {
							return err
						}

						if float64(n) != float64(int(n)) {
							return core.NotIntError(n)
						}

						w.WriteHeader(int(n))

						s, err := ts[0].EvalString()

						if err != nil {
							return err
						}

						if _, err := w.Write(([]byte)(s)); err != nil {
							return httpError(err)
						}

						return core.Nil
					})},
		})

	wg.Wait()
}
