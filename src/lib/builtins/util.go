package builtins

import "github.com/coel-lang/coel/src/lib/core"

func fileError(err error) core.Value {
	return core.NewError("FileError", err.Error())
}
