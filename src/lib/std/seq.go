package std

import "github.com/tisp-lang/tisp/src/lib/core"

// Seq runs arguments of outputs sequentially.
var Seq = core.NewLazyFunction(
	core.NewSignature(
		nil, nil, "outputs",
		nil, nil, "",
	),
	func(ts ...*core.Thunk) core.Value {
		l := ts[0]

		for {
			out := core.PApp(core.First, l)
			if err, ok := out.EvalOutput().(core.ErrorType); ok {
				return err
			}

			l = core.PApp(core.Rest, l)

			if v := checkEmptyList(l, out); v != nil {
				return v
			}
		}
	})
