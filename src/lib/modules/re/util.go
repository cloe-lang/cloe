package re

import "github.com/coel-lang/coel/src/lib/core"

func evaluateStringArguments(ts []*core.Thunk) ([]string, *core.Thunk) {
	ss := make([]string, 0, len(ts))

	for _, t := range ts {
		v := t.Eval()
		s, ok := v.(core.StringType)

		if !ok {
			return nil, core.NotStringError(v)
		}

		ss = append(ss, string(s))
	}

	return ss, nil
}
