package std

import "github.com/raviqqe/tisp/src/lib/core"

// List creates a list which contains elements of arguments.
var List = core.NewLazyFunction(
	core.NewSignature(
		nil, nil, "elems",
		nil, nil, "",
	),
	func(ts ...*core.Thunk) core.Value {
		return ts[0]
	})
