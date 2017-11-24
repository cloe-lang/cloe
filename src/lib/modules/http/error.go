package http

import "github.com/tisp-lang/tisp/src/lib/core"

func httpError(err error) *core.Thunk {
	return core.NewError("HTTPError", err.Error())
}
