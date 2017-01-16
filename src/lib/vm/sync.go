package vm

var Sync = NewLazyFunction(func(ts ...*Thunk) Object {
	if len(ts) != 2 {
		return numArgsError("sync", "2")
	}

	e, ok := ts[0].Eval().(errorType)

	if ok {
		return e
	}

	return ts[1]
})
