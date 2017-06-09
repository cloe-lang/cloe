package std

import "github.com/tisp-lang/tisp/src/lib/core"

// Y is Y combinator which takes a function whose first argument is itself
// applied to the combinator.
var Y = core.NewLazyFunction(
	core.NewSignature(
		[]string{"function"}, nil, "",
		nil, nil, "",
	),
	func(ts ...*core.Thunk) core.Value {
		if len(ts) != 1 {
			return core.NumArgsError("y", "1")
		}

		xfxx := core.PApp(core.Partial, fxx, ts[0])
		return core.PApp(xfxx, xfxx)
	})

var fxx = core.NewLazyFunction(
	core.NewSignature(
		[]string{"f", "x"}, nil, "",
		nil, nil, "",
	),
	func(ts ...*core.Thunk) core.Value {
		return core.PApp(core.Partial, ts[0], core.PApp(ts[1], ts[1]))
	})
