package http

import (
	"errors"
	"io/ioutil"
	"net/http"

	"github.com/tisp-lang/tisp/src/lib/core"
)

var get = core.NewLazyFunction(
	core.NewSignature(
		[]string{"url"}, nil, "",
		nil, []core.OptionalArgument{core.NewOptionalArgument("error", core.True)}, "",
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

		bs, err := ioutil.ReadAll(r.Body)

		if err != nil {
			return httpError(err)
		}

		v = ts[1].Eval()
		b, ok := v.(core.BoolType)

		if !ok {
			return core.NotBoolError(v)
		}

		if b && r.StatusCode/100 != 2 {
			return httpError(errors.New("status code is not 2XX"))
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
				core.NewString(string(bs)),
			})
	})
