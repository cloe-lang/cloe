package builtins

import (
	"github.com/tisp-lang/tisp/src/lib/core"
	"github.com/tisp-lang/tisp/src/lib/systemt"
)

// Par evaluates arguments asynchronously in parallel and returns the last one.
var Par = core.NewLazyFunction(
	core.NewSignature(
		nil, nil, "xs",
		nil, nil, "",
	),
	func(ts ...*core.Thunk) core.Value {
		l := ts[0]

		if v := checkEmptyList(l, core.NumArgsError("par", "> 0")); v != nil {
			return v
		}

		for {
			t := core.PApp(core.First, l)
			systemt.Daemonize(func() {
				t.Eval()
			})

			l = core.PApp(core.Rest, l)
			if v := checkEmptyList(l, t); v != nil {
				return v
			}
		}
	})
