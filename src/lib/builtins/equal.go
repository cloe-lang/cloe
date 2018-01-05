package builtins

import "github.com/coel-lang/coel/src/lib/core"

// Equal checks if all arguments are equal or not, and returns true if so or false otherwise.
var Equal = core.NewLazyFunction(
	core.NewSignature(nil, nil, "args", nil, nil, ""),
	func(ts ...*core.Thunk) core.Value {
		l := ts[0]

		if v := checkEmptyList(l, core.True); v != nil {
			return v
		}

		e := core.PApp(core.First, l)
		l = core.PApp(core.Rest, l)

		for {
			if v := checkEmptyList(l, core.True); v != nil {
				return v
			}

			v := core.PApp(core.Equal, e, core.PApp(core.First, l)).Eval()
			b, ok := v.(core.BoolType)

			if !ok {
				return core.NotBoolError(v)
			} else if !b {
				return core.False
			}

			l = core.PApp(core.Rest, l)
		}
	})
