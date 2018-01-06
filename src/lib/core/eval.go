package core

// EvalBool evaluates a thunk which is expected to be a boolean value.
func EvalBool(t *Thunk) (BoolType, Value) {
	v := t.Eval()
	b, ok := v.(BoolType)

	if !ok {
		return false, NotBoolError(v).Eval()
	}

	return b, nil
}

// EvalNumber evaluates a thunk which is expected to be a number value.
func EvalNumber(t *Thunk) (NumberType, Value) {
	v := t.Eval()
	n, ok := v.(NumberType)

	if !ok {
		return 0, NotNumberError(v).Eval()
	}

	return n, nil
}

// EvalString evaluates a thunk which is expected to be a string value.
func EvalString(t *Thunk) (StringType, Value) {
	v := t.Eval()
	s, ok := v.(StringType)

	if !ok {
		return "", NotStringError(v).Eval()
	}

	return s, nil
}
