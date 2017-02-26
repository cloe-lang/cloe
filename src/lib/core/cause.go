package core

var Cause = NewLazyFunction(
	NewSignature(
		[]string{"first", "next"}, []OptionalArgument{}, "",
		[]string{}, []OptionalArgument{}, "",
	),
	func(ts ...*Thunk) Object {
		e, ok := ts[0].Eval().(ErrorType)

		if ok {
			return e
		}

		return ts[1]
	})
