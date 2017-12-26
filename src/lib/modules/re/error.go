package re

import "github.com/coel-lang/coel/src/lib/core"

func regexError(err error) *core.Thunk {
	return core.NewError("RegexError", err.Error())
}
