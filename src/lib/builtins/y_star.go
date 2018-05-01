package builtins

import "github.com/cloe-lang/cloe/src/lib/core"

// Ys is Y* combinator which takes functions whose first arguments are a list
// of themselves applied to the combinator.
var Ys = core.NewLazyFunction(
	core.NewSignature(nil, "functions", nil, ""),
	func(ts ...core.Value) core.Value {
		t := ts[0]

		return core.PApp(xx, core.NewLazyFunction(
			core.NewSignature([]string{"x"}, "", nil, ""),
			func(ts ...core.Value) core.Value {
				s := ts[0]

				applyF := core.NewLazyFunction(
					core.NewSignature([]string{"f"}, "args", nil, "kwargs"),
					func(ts ...core.Value) core.Value {
						return core.App(ts[0], core.NewArguments(
							[]core.PositionalArgument{
								core.NewPositionalArgument(core.PApp(s, s), false),
								core.NewPositionalArgument(ts[1], true),
							},
							[]core.KeywordArgument{core.NewKeywordArgument("", ts[2])}))
					})

				return createNewFuncs(t, applyF)
			}))
	})

func createNewFuncs(olds, applyF core.Value) core.Value {
	if v := core.ReturnIfEmptyList(olds, core.EmptyList); v != nil {
		return v
	}

	return core.StrictPrepend(
		[]core.Value{core.PApp(core.Partial, applyF, core.PApp(core.First, olds))},
		createNewFuncs(core.PApp(core.Rest, olds), applyF))
}

var xx = core.NewLazyFunction(
	core.NewSignature([]string{"x"}, "", nil, ""),
	func(ts ...core.Value) core.Value {
		return core.PApp(ts[0], ts[0])
	})
