package std

import "../core"

var Cause = core.NewLazyFunction(
	core.NewSignature(
		[]string{"first", "next"}, []core.OptionalArgument{}, "",
		[]string{}, []core.OptionalArgument{}, "",
	),
	func(ts ...*core.Thunk) core.Object {
		e, ok := ts[0].Eval().(core.ErrorType)

		if ok {
			return e
		}

		return ts[1]
	})
