package re

import "github.com/cloe-lang/cloe/src/lib/core"

func regexError(err error) core.Value {
	return core.NewError("RegexError", err.Error())
}
