package re

import "github.com/coel-lang/coel/src/lib/core"

func regexError(err error) core.Value {
	return core.NewError("RegexError", err.Error())
}
