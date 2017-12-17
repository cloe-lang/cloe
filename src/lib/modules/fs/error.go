package fs

import "github.com/coel-lang/coel/src/lib/core"

func fileSystemError(err error) *core.Thunk {
	return core.NewError("FileSystemError", err.Error())
}
