package builtins

import (
	"github.com/coel-lang/coel/src/lib/core"
	"github.com/coel-lang/coel/src/lib/systemt"
)

// Par evaluates arguments in parallel and returns the last one.
var Par = core.NewLazyFunction(
	core.NewSignature(nil, nil, "args", nil, nil, ""),
	func(ts ...core.Value) core.Value {
		l := ts[0]

		for {
			t := core.PApp(core.First, l)
			l = core.PApp(core.Rest, l)

			if v := core.ReturnIfEmptyList(l, t); v != nil {
				return v
			}

			systemt.Daemonize(func() { t.Eval() })
		}
	})
