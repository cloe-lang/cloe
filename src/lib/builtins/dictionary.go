package builtins

import "github.com/tisp-lang/tisp/src/lib/core"

// Dictionary creates a new dictionary from pairs of a key and value.
var Dictionary = core.NewLazyFunction(
	core.NewSignature(
		nil, nil, "keyValuePairs",
		nil, nil, "",
	),
	func(ts ...*core.Thunk) core.Value {
		d := core.EmptyDictionary
		l := ts[0]

		for {
			if v := checkEmptyList(l, d); v != nil {
				return v
			}

			k := core.PApp(core.First, l)
			l = core.PApp(core.Rest, l)
			val := core.PApp(core.First, l)
			l = core.PApp(core.Rest, l)

			d = core.PApp(core.Insert, d, k, val)
		}
	})
