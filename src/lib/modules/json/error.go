package json

import "github.com/coel-lang/coel/src/lib/core"

func jsonError(err error) *core.Thunk {
	return core.NewError("JSONError", err.Error())
}
