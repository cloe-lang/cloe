package core

// IsOrdered checks if a value is ordered or not.
var IsOrdered = NewLazyFunction(
	NewSignature([]string{"arg"}, nil, "", nil, nil, ""),
	isOrderedFunction)

func isOrderedFunction(ts ...*Thunk) Value {
	switch x := ts[0].Eval().(type) {
	case ErrorType:
		return x
	case ListType:
		for !x.Empty() {
			v := ensureNormal(isOrderedFunction(x.First()))
			b, ok := v.(BoolType)

			if !ok {
				return NotBoolError(v)
			} else if !b {
				return False
			}

			var err Value
			if x, err = x.Rest().EvalList(); err != nil {
				return err
			}
		}

		return True
	default:
		_, ok := x.(ordered)
		return NewBool(ok)
	}
}
