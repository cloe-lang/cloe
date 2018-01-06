package core

// IsOrdered checks if a value is ordered or not.
var IsOrdered = NewLazyFunction(
	NewSignature([]string{"arg"}, nil, "", nil, nil, ""),
	isOrderedFunction)

func isOrderedFunction(ts ...*Thunk) Value {
	v := ts[0].Eval()

	if e, ok := v.(ErrorType); ok {
		return e
	}

	if _, ok := v.(ListType); !ok {
		_, ok := v.(ordered)
		return NewBool(ok)
	}

	l := ts[0]

	for {
		if v := ReturnIfEmptyList(l, True); v != nil {
			return v
		}

		t := PApp(First, l)
		l = PApp(Rest, l)

		v := ensureNormal(isOrderedFunction(t))
		b, ok := v.(BoolType)

		if !ok {
			return NotBoolError(v)
		} else if !b {
			return False
		}
	}
}
