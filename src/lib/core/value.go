package core

// Value represents a value.
type Value interface {
	// Eval svaluates a value into WHNF. (i.e. Thunks would be unwrapped.)
	eval() Value
}

// EvalBoolean evaluates a thunk which is expected to be a boolean value.
func EvalBoolean(v Value) (BooleanType, Value) {
	b, ok := EvalPure(v).(*BooleanType)

	if !ok {
		return *False, NotBooleanError(v)
	}

	return *b, nil
}

// EvalDictionary evaluates a thunk which is expected to be a dictionary value.
func EvalDictionary(v Value) (*DictionaryType, Value) {
	d, ok := EvalPure(v).(*DictionaryType)

	if !ok {
		return EmptyDictionary, NotDictionaryError(v)
	}

	return d, nil
}

// EvalList evaluates a thunk which is expected to be a list value.
func EvalList(v Value) (*ListType, Value) {
	l, ok := EvalPure(v).(*ListType)

	if !ok {
		return EmptyList, NotListError(v)
	}

	return l, nil
}

// EvalNumber evaluates a thunk which is expected to be a number value.
func EvalNumber(v Value) (NumberType, Value) {
	n, ok := EvalPure(v).(*NumberType)

	if !ok {
		return 0, NotNumberError(v)
	}

	return *n, nil
}

// EvalString evaluates a thunk which is expected to be a string value.
func EvalString(v Value) (StringType, Value) {
	s, ok := EvalPure(v).(StringType)

	if !ok {
		return "", NotStringError(v)
	}

	return s, nil
}

func evalCollection(v Value) (collection, Value) {
	c, ok := EvalPure(v).(collection)

	if !ok {
		return nil, NotCollectionError(v)
	}

	return c, nil
}

func evalSequence(v Value) (sequence, Value) {
	s, ok := EvalPure(v).(sequence)

	if !ok {
		return nil, NotSequenceError(v)
	}

	return s, nil
}

// EvalPure evaluates a pure value.
func EvalPure(v Value) Value {
	v = v.eval()
	_, ok := v.(effectType)

	if ok {
		return impureFunctionError()
	}

	return v
}

// EvalImpure evaluates an impure function call.
func EvalImpure(v Value) Value {
	e, ok := v.eval().(effectType)

	if !ok {
		return NotEffectError(v)
	}

	return EvalPure(e.value)
}
