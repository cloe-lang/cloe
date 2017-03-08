package std

import "github.com/raviqqe/tisp/src/lib/core"

// Cause runs arguments of outputs sequentially.
var Cause = core.NewLazyFunction(
	core.NewSignature(
		[]string{}, []core.OptionalArgument{}, "outputs",
		[]string{}, []core.OptionalArgument{}, "",
	),
	func(ts ...*core.Thunk) core.Object {
		t := ts[0]
		out := core.PApp(core.First, t)

		for {
			if err, ok := out.Eval().(core.ErrorType); ok {
				return err
			}

			t = core.PApp(core.Rest, t)

			o := core.PApp(core.Equal, t, core.EmptyList).Eval()
			b, ok := o.(core.BoolType)

			if !ok {
				return core.NotBoolError(o)
			} else if b {
				break
			}

			out = core.PApp(core.First, t)
		}

		return out
	})
