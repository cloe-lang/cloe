package http

import "github.com/coel-lang/coel/src/lib/core"

func httpError(err error) core.Value {
	return core.NewError("HTTPError", err.Error())
}
