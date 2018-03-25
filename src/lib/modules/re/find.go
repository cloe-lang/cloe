package re

import (
	"regexp"

	"github.com/coel-lang/coel/src/lib/core"
)

var find = core.NewLazyFunction(
	core.NewSignature([]string{"pattern", "src"}, "", nil, ""),
	func(vs ...core.Value) core.Value {
		ss, e := evaluateStringArguments(vs)

		if e != nil {
			return e
		}

		r, err := regexp.Compile(ss[0])

		if err != nil {
			return regexError(err)
		}

		ss = r.FindStringSubmatch(ss[1])

		if len(ss) == 0 {
			return core.Nil
		}

		vs = make([]core.Value, 0, len(ss))

		for _, s := range ss {
			v := core.Value(core.Nil)

			if s != "" {
				v = core.NewString(s)
			}

			vs = append(vs, v)
		}

		return core.NewList(vs...)
	})
