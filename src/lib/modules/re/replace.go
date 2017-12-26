package re

import (
	"regexp"

	"github.com/coel-lang/coel/src/lib/core"
)

var replace = core.NewLazyFunction(
	core.NewSignature([]string{"pattern", "repl", "src"}, nil, "", nil, nil, ""),
	func(ts ...*core.Thunk) core.Value {
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
