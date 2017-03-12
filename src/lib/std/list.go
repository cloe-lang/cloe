package std

import "github.com/raviqqe/tisp/src/lib/core"

// List creates a list which contains elements of arguments.
var List = core.NewLazyFunction(
	core.NewSignature(
		[]string{}, []core.OptionalArgument{}, "elems",
		[]string{}, []core.OptionalArgument{}, "",
	),
	func(ts ...*core.Thunk) core.Value {
		return ts[0]
	})
