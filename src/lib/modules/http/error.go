package http

import "github.com/coel-lang/coel/src/lib/core"

func httpError(err error) *core.Thunk {
	return core.NewError("HTTPError", err.Error())
}
