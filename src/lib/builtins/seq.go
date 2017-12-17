package builtins

import (
	"github.com/coel-lang/coel/src/lib/core"
)

// Seq runs arguments of effects sequentially and returns the last one.
var Seq = core.NewLazyFunction(
	core.NewSignature(nil, nil, "effects", nil, nil, ""),
	func(ts ...*core.Thunk) core.Value {
		l := ts[0]

		for {
			t := core.PApp(core.First, l)

			if err, ok := t.EvalEffect().(core.ErrorType); ok {
				return err
			}

			l = core.PApp(core.Rest, l)

			if v := checkEmptyList(l, t); v != nil {
				return v
			}
		}
	})
