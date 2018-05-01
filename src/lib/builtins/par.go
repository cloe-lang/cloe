package builtins

import (
	"github.com/cloe-lang/cloe/src/lib/core"
	"github.com/cloe-lang/cloe/src/lib/systemt"
)

// Par evaluates arguments in parallel and returns the last one.
var Par = core.NewLazyFunction(
	core.NewSignature(nil, "args", nil, ""),
	func(ts ...core.Value) core.Value {
		l := ts[0]

		for {
			v := core.PApp(core.First, l)
			l = core.PApp(core.Rest, l)

			if v := core.ReturnIfEmptyList(l, v); v != nil {
				return v
			}

			systemt.Daemonize(func() { core.EvalPure(v) })
		}
	})
