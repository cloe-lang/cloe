package http

import (
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
			core.NewOptionalArgument("error", core.True),
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

		return handleMethodResult(r, err, ts[3])
	})
