package std

import "github.com/tisp-lang/tisp/src/lib/core"

// Ys is Y* combinator which takes functions whose first arguments are a list
// of themselves applied to the combinator.
var Ys = core.NewLazyFunction(
	core.NewSignature(
		nil, nil, "functions",
		nil, nil, "",
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
				[]string{"x"}, nil, "",
				nil, nil, "",
			),
			func(ps ...*core.Thunk) core.Value {
				p := ps[0]

				applyF := core.NewLazyFunction(
					core.NewSignature(
						[]string{"f"}, nil, "args",
						nil, nil, "kwargs",
					),
					func(qs ...*core.Thunk) core.Value {
						return core.App(qs[0], core.NewArguments(
							[]core.PositionalArgument{
								core.NewPositionalArgument(core.PApp(p, p), false),
								core.NewPositionalArgument(qs[1], true),
							},
							nil,
							[]*core.Thunk{qs[2]}))
					})

				newFs := make([]*core.Thunk, 0, len(fs))

				for _, f := range fs {
					newFs = append(newFs, core.PApp(core.Partial, applyF, f))
				}

				return core.NewList(newFs...)
			})

		return core.PApp(xx, f)
	})

var xx = core.NewLazyFunction(
	core.NewSignature(
		[]string{"x"}, nil, "",
		nil, nil, "",
	),
	func(ts ...*core.Thunk) core.Value {
		return core.PApp(ts[0], ts[0])
	})
