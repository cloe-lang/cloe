package re

import "github.com/coel-lang/coel/src/lib/core"

func evaluateStringArguments(vs []core.Value) ([]string, core.Value) {
	ss := make([]string, 0, len(vs))

	for _, v := range vs {
		s, err := core.EvalString(v)

		if err != nil {
			return nil, err
		}

		ss = append(ss, string(s))
	}

	return ss, nil
}
