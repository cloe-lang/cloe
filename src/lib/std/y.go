package std

import "github.com/tisp-lang/tisp/src/lib/core"

// Y is Y combinator which takes a function whose first argument is itself
// applied to the combinator.
// Using Y combinator to define built-in functions in Go source is dangerous
// because top-level recursive functions generate infinitely nested closures.
// (i.e. closure{f, x} where x will also be evaluated as closure{f, x}.)
var Y = core.NewLazyFunction(
	core.NewSignature(
		[]string{"function"}, nil, "",
		nil, nil, "",
	),
	func(ts ...*core.Thunk) core.Value {
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
