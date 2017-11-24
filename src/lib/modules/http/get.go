package http

import (
	"io/ioutil"
	"net/http"

	"github.com/tisp-lang/tisp/src/lib/core"
)

var get = core.NewLazyFunction(
	core.NewSignature(
		[]string{"url"}, nil, "",
		nil, nil, "",
	),
	func(ts ...*core.Thunk) core.Value {
		v := ts[0].Eval()
		s, ok := v.(core.StringType)

		if !ok {
			return core.NotStringError(v)
		}

		r, err := http.Get(string(s))

		if err != nil {
			return httpError(err)
		}

		b, err := ioutil.ReadAll(r.Body)

		if err != nil {
			return httpError(err)
		}

		if err = r.Body.Close(); err != nil {
			return httpError(err)
		}

		return core.NewDictionary(
			[]core.Value{
				core.NewString("status").Eval(),
				core.NewString("body").Eval(),
			},
			[]*core.Thunk{
				core.NewNumber(float64(r.StatusCode)),
				core.NewString(string(b)),
			})
	})
