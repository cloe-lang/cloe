package os

import "github.com/cloe-lang/cloe/src/lib/core"

func osError(err error) core.Value {
	return core.NewError("OSError", err.Error())
}
