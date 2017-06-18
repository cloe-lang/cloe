package std

import "github.com/tisp-lang/tisp/src/lib/core"

// TODO: Implement LessEq, Greater, and GreaterEq functions.

func compare(checkOrder func(core.NumberType) bool) func(ts ...*core.Thunk) core.Value {
	return func(ts ...*core.Thunk) core.Value {
		l := ts[0]

		if v := checkEmptyList(l, core.True); v != nil {
			return v
		}

		prev := core.PApp(core.First, l)

		for {
			l = core.PApp(core.Rest, l)

			if v := checkEmptyList(l, core.True); v != nil {
				return v
			}

			current := core.PApp(core.First, l)

			v := core.PApp(core.Compare, prev, current).Eval()
			n, ok := v.(core.NumberType)

			if !ok {
				return core.NotBoolError(v)
			} else if !checkOrder(n) {
				return core.False
			}

			prev = current
		}
	}
}

// Less checks if arguments are aligned in ascending order or not.
var Less = core.NewLazyFunction(
	core.NewSignature([]string{}, nil, "xs", nil, nil, ""),
	compare(func(n core.NumberType) bool { return n == -1 }))
