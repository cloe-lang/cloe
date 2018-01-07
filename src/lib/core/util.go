package core

import (
	"fmt"
)

func sprint(x interface{}) StringType {
	return StringType(fmt.Sprint(x))
}

type dumpable interface {
	dump() Value
}

// Dump dumps a value into a string type value.
var Dump = NewLazyFunction(
	NewSignature([]string{"arg"}, nil, "", nil, nil, ""),
	func(ts ...*Thunk) Value {
		s, err := StrictDump(ts[0].Eval())

		if err != nil {
			return err
		}

		return s
	})

// StrictDump is a variant of Dump which evaluates input strictly.
func StrictDump(v Value) (StringType, Value) {
	switch x := ensureNormal(v).(type) {
	case ErrorType:
		return "", x
	case dumpable:
		v = x.dump()
	case stringable:
		v = x.string()
	default:
		panic(fmt.Errorf("Invalid value: %#v", x))
	}

	s, ok := ensureNormal(v).(StringType)

	if !ok {
		return "", NotStringError(v).Eval()
	}

	return s, nil
}

// ensureNormal evaluates nested thunks into WHNF values.
// This function must be used with care because it prevents tail call
// elimination.
func ensureNormal(v Value) Value {
	if t, ok := v.(*Thunk); ok {
		return t.Eval()
	}

	return v
}

var identity = NewLazyFunction(
	NewSignature([]string{"arg"}, nil, "", nil, nil, ""),
	func(ts ...*Thunk) Value { return ts[0] })

// TypeOf returns a type name of an argument as a string.
var TypeOf = NewLazyFunction(
	NewSignature([]string{"arg"}, nil, "", nil, nil, ""),
	func(ts ...*Thunk) Value {
		// No case of effectType should be here.
		switch ts[0].Eval().(type) {
		case BoolType:
			return NewString("bool")
		case DictionaryType:
			return NewString("dict")
		case ListType:
			return NewString("list")
		case NilType:
			return NewString("nil")
		case NumberType:
			return NewString("number")
		case StringType:
			return NewString("string")

		case functionType:
			return NewString("function")
		case RawFunctionType:
			return NewString("function")

		case ErrorType:
			// TODO: Remove this case and use catch function to check if a value is an
			// error or not.
			return NewString("error")
		}

		panic(fmt.Errorf("Invalid value: %#v", ts[0].Eval()))
	})
