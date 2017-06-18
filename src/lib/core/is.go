package core

// IsOrdered checks if a value is ordered or not.
var IsOrdered = NewLazyFunction(
	NewSignature([]string{"x"}, nil, "", nil, nil, ""),
	func(ts ...*Thunk) Value {
		_, ok := ts[0].Eval().(ordered)
		return NewBool(ok)
	})
