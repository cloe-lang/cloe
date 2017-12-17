package builtins

import (
	"github.com/coel-lang/coel/src/lib/core"
	"github.com/coel-lang/coel/src/lib/systemt"
)

// Par evaluates arguments in parallel and returns the last one.
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
