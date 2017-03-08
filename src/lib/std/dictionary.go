package std

import "github.com/raviqqe/tisp/src/lib/core"

// Dictionary creates a new dictionary from pairs of a key and value.
var Dictionary = core.NewLazyFunction(
	core.NewSignature(
		[]string{}, []core.OptionalArgument{}, "keyValuePairs",
		[]string{}, []core.OptionalArgument{}, "",
	),
	func(ts ...*core.Thunk) core.Object {
		d := core.EmptyDictionary
		l := ts[0]

		for {
			o := core.PApp(core.Equal, l, core.EmptyList).Eval()
			b, ok := o.(core.BoolType)

			if !ok {
				return core.NotBoolError(o)
			} else if b {
				break
			}

			k := core.PApp(core.First, l)
			l = core.PApp(core.Rest, l)
			v := core.PApp(core.First, l)
			l = core.PApp(core.Rest, l)

			d = core.PApp(core.Set, d, k, v)
		}

		return d
	})
