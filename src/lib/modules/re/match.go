package re

import (
	"regexp"

	"github.com/cloe-lang/cloe/src/lib/core"
)

var match = core.NewLazyFunction(
	core.NewSignature([]string{"pattern", "src"}, "", nil, ""),
	func(ts ...core.Value) core.Value {
		ss, e := evaluateStringArguments(ts)

		if e != nil {
			return e
		}

		b, err := regexp.MatchString(ss[0], ss[1])

		if err != nil {
			return regexError(err)
		}

		return core.NewBoolean(b)
	})
