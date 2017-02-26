package core

var Y = NewLazyFunction(
	NewSignature(
		[]string{"function"}, []OptionalArgument{}, "",
		[]string{}, []OptionalArgument{}, "",
	),
	func(ts ...*Thunk) Object {
		if len(ts) != 1 {
			return NumArgsError("y", "1")
		}

		xfxx := PApp(Partial, fxx, ts[0])
		return PApp(xfxx, xfxx)
	})

var fxx = NewLazyFunction(
	NewSignature(
		[]string{"f", "x"}, []OptionalArgument{}, "",
		[]string{}, []OptionalArgument{}, "",
	),
	func(ts ...*Thunk) Object {
		return PApp(Partial, ts[0], PApp(ts[1], ts[1]))
	})

var Ys = NewLazyFunction(
	NewSignature(
		[]string{}, []OptionalArgument{}, "functions",
		[]string{}, []OptionalArgument{}, "",
	),
	func(ts ...*Thunk) Object {
		o := ts[0].Eval()
		l, ok := o.(ListType)

		if !ok {
			panic("Rest arguments must be a list.")
		}

		fs, err := l.ToThunks()

		if err != nil {
			return err
		}

		if len(fs) == 0 {
			return NumArgsError("ys", ">= 1")
		}

		f := NewLazyFunction(
			NewSignature(
				[]string{"x"}, []OptionalArgument{}, "",
				[]string{}, []OptionalArgument{}, "",
			),
			func(ps ...*Thunk) Object {
				p := ps[0]

				applyF := NewLazyFunction(
					NewSignature(
						[]string{"f"}, []OptionalArgument{}, "args",
						[]string{}, []OptionalArgument{}, "kwargs",
					),
					func(qs ...*Thunk) Object {
						return App(qs[0], NewArguments(
							[]PositionalArgument{
								NewPositionalArgument(PApp(p, p), false),
								NewPositionalArgument(qs[1], true),
							},
							[]KeywordArgument{},
							[]*Thunk{qs[2]}))
					})

				newFs := make([]*Thunk, len(fs))

				for i, f := range fs {
					newFs[i] = PApp(Partial, applyF, f)
				}

				return NewList(newFs...)
			})

		return PApp(xx, f)
	})

var xx = NewLazyFunction(
	NewSignature(
		[]string{"x"}, []OptionalArgument{}, "",
		[]string{}, []OptionalArgument{}, "",
	),
	func(ts ...*Thunk) Object {
		return PApp(ts[0], ts[0])
	})
