package http

import "github.com/cloe-lang/cloe/src/lib/core"

func httpError(err error) core.Value {
	return core.NewError("HTTPError", err.Error())
}
