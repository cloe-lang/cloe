package re

import (
	"regexp"

	"github.com/coel-lang/coel/src/lib/core"
)

var find = core.NewLazyFunction(
	core.NewSignature([]string{"pattern", "src"}, nil, "", nil, nil, ""),
	func(ts ...*core.Thunk) core.Value {
		ss, e := evaluateStringArguments(ts)

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

		ts = make([]*core.Thunk, 0, len(ss))

		for _, s := range ss {
			t := core.Nil

			if s != "" {
				t = core.NewString(s)
			}

			ts = append(ts, t)
		}

		return core.NewList(ts...)
	})
