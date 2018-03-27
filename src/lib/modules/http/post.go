package http

import (
	"net/http"
	"strings"

	"github.com/coel-lang/coel/src/lib/core"
)

var post = core.NewLazyFunction(
	core.NewSignature(
		[]string{"url", "body"}, "",
		[]core.OptionalParameter{
			core.NewOptionalParameter("contentType", core.NewString("text/plain")),
			core.NewOptionalParameter("error", core.True),
		},
		"",
	),
	func(vs ...core.Value) core.Value {
		ss := make([]string, 0, 3)

		for i := 0; i < cap(ss); i++ {
			s, err := core.EvalString(vs[i])

			if err != nil {
				return err
			}

			ss = append(ss, string(*s))
		}

		r, err := http.Post(ss[0], ss[2], strings.NewReader(ss[1]))

		return handleMethodResult(r, err, vs[3])
	})
