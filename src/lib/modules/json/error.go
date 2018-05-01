package json

import "github.com/cloe-lang/cloe/src/lib/core"

func jsonError(err error) core.Value {
	return core.NewError("JSONError", err.Error())
}
