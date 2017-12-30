package core

var isOrderedSignature = NewSignature([]string{"arg"}, nil, "", nil, nil, "")

// IsOrdered checks if a value is ordered or not.
var IsOrdered = NewLazyFunction(isOrderedSignature, isOrderedFunction)

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
		v := PApp(Equal, l, EmptyList).Eval()
		b, ok := v.(BoolType)

		if !ok {
			return NotBoolError(v)
		} else if b {
			return True
		}

		t := PApp(First, l)
		l = PApp(Rest, l)

		v = PApp(NewLazyFunction(isOrderedSignature, isOrderedFunction), t).Eval()
		b, ok = v.(BoolType)

		if !ok {
			return NotBoolError(v)
		} else if !b {
			return False
		}
	}
}
