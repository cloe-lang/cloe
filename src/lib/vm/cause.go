package vm

var Cause = NewLazyFunction(func(ts ...*Thunk) Object {
	if len(ts) != 2 {
		return NumArgsError("cause", "2")
	}

	e, ok := ts[0].Eval().(ErrorType)

	if ok {
		return e
	}

	return ts[1]
})
