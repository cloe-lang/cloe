package http

import (
	"net/http"
	"sync"

	"github.com/tisp-lang/tisp/src/lib/core"
	"github.com/tisp-lang/tisp/src/lib/std"
	"github.com/tisp-lang/tisp/src/lib/systemt"
)

const requestChannelSize = 1024
const responseChannelSize = 1024

var getRequests = core.NewLazyFunction(
	core.NewSignature(
		[]string{"address"}, nil, "",
		nil, nil, "",
	),
	func(ts ...*core.Thunk) core.Value {
		v := ts[0].Eval()
		s, ok := v.(core.StringType)

		if !ok {
			return core.NotStringError(v)
		}

		ec := make(chan error)
		h := newHandler()

		systemt.Daemonize(func() {
			if err := http.ListenAndServe(string(s), h); err != nil {
				ec <- err
			}
		})

		return core.PApp(core.PApp(std.Y, core.NewLazyFunction(
			core.NewSignature([]string{"me"}, nil, "", nil, nil, ""),
			func(ts ...*core.Thunk) core.Value {
				select {
				case t := <-h.Requests:
					return core.PApp(core.Prepend, t, core.PApp(ts[0]))
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
	wg := sync.WaitGroup{}

	wg.Add(1)

	h.Requests <- core.NewDictionary(
		[]core.Value{
			core.NewString("respond").Eval(),
		},
		[]*core.Thunk{
			core.NewLazyFunction(
				core.NewSignature([]string{"body"}, nil, "", nil, nil, ""),
				func(ts ...*core.Thunk) core.Value {
					defer wg.Done()

					v := ts[0].Eval()
					s, ok := v.(core.StringType)

					if !ok {
						return core.NotStringError(v)
					}

					if _, err := w.Write(([]byte)(s)); err != nil {
						return httpError(err)
					}

					return core.NewEffect(core.Nil)
				}),
		})

	wg.Wait()
}
