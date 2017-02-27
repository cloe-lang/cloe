package std

import "github.com/raviqqe/tisp/src/lib/core"

var List = core.NewLazyFunction(
	core.NewSignature(
		[]string{}, []core.OptionalArgument{}, "elems",
		[]string{}, []core.OptionalArgument{}, "",
	),
	func(ts ...*core.Thunk) core.Object {
		return ts[0]
	})
