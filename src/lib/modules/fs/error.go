package fs

import "github.com/cloe-lang/cloe/src/lib/core"

func fileSystemError(err error) core.Value {
	return core.NewError("FileSystemError", err.Error())
}
