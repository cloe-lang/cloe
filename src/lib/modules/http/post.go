package http

import (
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/tisp-lang/tisp/src/lib/core"
)

var post = core.NewLazyFunction(
	core.NewSignature(
		[]string{"url", "body"}, nil, "",
		nil,
		[]core.OptionalArgument{
			core.NewOptionalArgument("contentType", core.NewString("text/plain")),
		},
		"",
	),
	func(ts ...*core.Thunk) core.Value {
		ss := make([]string, 0, 3)

		for i := 0; i < cap(ss); i++ {
			v := ts[i].Eval()
			s, ok := v.(core.StringType)

			if !ok {
				return core.NotStringError(v)
			}

			ss = append(ss, string(s))
		}

		r, err := http.Post(ss[0], ss[2], strings.NewReader(ss[1]))

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
