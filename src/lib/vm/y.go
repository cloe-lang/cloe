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
