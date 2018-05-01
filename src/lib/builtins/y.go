package builtins

import "github.com/cloe-lang/cloe/src/lib/core"

// Y is Y combinator which takes a function whose first argument is itself
// applied to the combinator.
//
// THE COMMENT BELOW MAY BE OUTDATED because we moved from a lambda calculus
// based combinator to an implementation based on a recursive function in Go.
//
// Using Y combinator to define built-in functions in Go source is dangerous
// because top-level recursive functions generate infinitely nested closures.
// (i.e. closure{f, x} where x will also be evaluated as closure{f, x}.)
var Y = core.NewLazyFunction(
	core.NewSignature([]string{"function"}, "", nil, ""),
	func(ts ...core.Value) core.Value {
		return y(ts[0])
	})

func y(f core.Value) core.Value {
	return core.NewRawFunction(func(args core.Arguments) core.Value {
		return core.App(f, core.NewPositionalArguments(y(f)).Merge(args))
	})
}
