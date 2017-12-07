package builtins

import (
	"github.com/tisp-lang/tisp/src/lib/core"
)

// Seq runs arguments of effects sequentially.
var Seq = core.NewLazyFunction(
	core.NewSignature(
		nil, nil, "effects",
		nil, nil, "",
	),
	func(ts ...*core.Thunk) core.Value {
		l := ts[0]

		for {
			out := core.PApp(core.First, l)
			if err, ok := out.EvalEffect().(core.ErrorType); ok {
				return err
			}

			l = core.PApp(core.Rest, l)

			if v := checkEmptyList(l, out); v != nil {
				return v
			}
		}
	})
