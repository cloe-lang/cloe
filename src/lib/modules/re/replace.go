package re

import (
	"regexp"

	"github.com/cloe-lang/cloe/src/lib/core"
)

var replace = core.NewLazyFunction(
	core.NewSignature([]string{"pattern", "repl", "src"}, "", nil, ""),
	func(ts ...core.Value) core.Value {
		ss, e := evaluateStringArguments(ts)

		if e != nil {
			return e
		}

		r, err := regexp.Compile(ss[0])

		if err != nil {
			return regexError(err)
		}

		return core.NewString(r.ReplaceAllString(ss[2], ss[1]))
	})
