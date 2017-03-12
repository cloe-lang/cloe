package std

import "github.com/raviqqe/tisp/src/lib/core"

var Y = core.NewLazyFunction(
	core.NewSignature(
		[]string{"function"}, []core.OptionalArgument{}, "",
		[]string{}, []core.OptionalArgument{}, "",
	),
	func(ts ...*core.Thunk) core.Value {
		if len(ts) != 1 {
			return core.NumArgsError("y", "1")
		}

		xfxx := core.PApp(core.Partial, fxx, ts[0])
		return core.PApp(xfxx, xfxx)
	})

var fxx = core.NewLazyFunction(
	core.NewSignature(
		[]string{"f", "x"}, []core.OptionalArgument{}, "",
		[]string{}, []core.OptionalArgument{}, "",
	),
	func(ts ...*core.Thunk) core.Value {
		return core.PApp(core.Partial, ts[0], core.PApp(ts[1], ts[1]))
	})

var Ys = core.NewLazyFunction(
	core.NewSignature(
		[]string{}, []core.OptionalArgument{}, "functions",
		[]string{}, []core.OptionalArgument{}, "",
	),
	func(ts ...*core.Thunk) core.Value {
		v := ts[0].Eval()
		l, ok := v.(core.ListType)

		if !ok {
			panic("Rest arguments must be a list.")
		}

		fs, err := l.ToThunks()

		if err != nil {
			return err
		}

		if len(fs) == 0 {
			return core.NumArgsError("ys", ">= 1")
		}

		f := core.NewLazyFunction(
			core.NewSignature(
				[]string{"x"}, []core.OptionalArgument{}, "",
				[]string{}, []core.OptionalArgument{}, "",
			),
			func(ps ...*core.Thunk) core.Value {
				p := ps[0]

				applyF := core.NewLazyFunction(
					core.NewSignature(
						[]string{"f"}, []core.OptionalArgument{}, "args",
						[]string{}, []core.OptionalArgument{}, "kwargs",
					),
					func(qs ...*core.Thunk) core.Value {
						return core.App(qs[0], core.NewArguments(
							[]core.PositionalArgument{
								core.NewPositionalArgument(core.PApp(p, p), false),
								core.NewPositionalArgument(qs[1], true),
							},
							[]core.KeywordArgument{},
							[]*core.Thunk{qs[2]}))
					})

				newFs := make([]*core.Thunk, len(fs))

				for i, f := range fs {
					newFs[i] = core.PApp(core.Partial, applyF, f)
				}

				return core.NewList(newFs...)
			})

		return core.PApp(xx, f)
	})

var xx = core.NewLazyFunction(
	core.NewSignature(
		[]string{"x"}, []core.OptionalArgument{}, "",
		[]string{}, []core.OptionalArgument{}, "",
	),
	func(ts ...*core.Thunk) core.Value {
		return core.PApp(ts[0], ts[0])
	})
