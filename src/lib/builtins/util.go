package builtins

import "github.com/cloe-lang/cloe/src/lib/core"

func fileError(err error) core.Value {
	return core.NewError("FileError", err.Error())
}
