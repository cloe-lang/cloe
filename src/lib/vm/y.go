package vm

var Y = NewLazyFunction(y)

func y(ts ...*Thunk) Object {
	if len(ts) != 1 {
		return numArgsError("y", "1")
	}

	xfxx := App(Partial, fxx, ts[0])
	return App(xfxx, xfxx)
}

var fxx = NewLazyFunction(fxxImpl)

func fxxImpl(ts ...*Thunk) Object {
	return App(Partial, ts[0], App(ts[1], ts[1]))
}

var Ys = NewLazyFunction(ys) // TODO: Test Ys more.

func ys(fs ...*Thunk) Object {
	// Note that ys returns []*Thunk. Results should not be applied to other
	// functions directly. You may want to wrap them using NewList.

	if len(fs) == 0 {
		return numArgsError("ys", ">= 1")
	}

	f := NewLazyFunction(func(ps ...*Thunk) Object {
		if len(ps) != 1 {
			panic("f takes only a argument.")
		}

		p := ps[0]

		applyFs := NewLazyFunction(func(qs ...*Thunk) Object {
			return App(qs[0], append(App(p, p).Eval().([]*Thunk), qs[1:]...)...)
		})

		newFs := make([]*Thunk, len(fs))

		for i, f := range fs {
			newFs[i] = App(Partial, applyFs, f)
		}

		return newFs
	})

	return App(xx, f)
}

var xx = NewLazyFunction(xxImpl)

func xxImpl(ts ...*Thunk) Object {
	if len(ts) != 1 {
		panic("xx takes only one argument.")
	}

	return App(ts[0], ts[0])
}
