package core

import (
	"fmt"
)

func sprint(s interface{}) StringType {
	return NewString(fmt.Sprint(s))
}

// Dump dumps a value into a string type value.
var Dump = NewLazyFunction(
	NewSignature([]string{"arg"}, nil, "", nil, nil, ""),
	func(vs ...Value) Value {
		s, err := StrictDump(vs[0])

		if err != nil {
			return err
		}

		return s
	})

// StrictDump is a variant of Dump which evaluates input strictly.
func StrictDump(v Value) (StringType, Value) {
	switch x := EvalPure(v).(type) {
	case ErrorType:
		return "", x
	case StringType:
		v = x.quoted()
	case stringable:
		v = x.string()
	default:
		panic(fmt.Errorf("Invalid value: %#v", x))
	}

	s, err := EvalString(v)

	if err != nil {
		return "", err
	}

	return s, nil
}

var identity = NewLazyFunction(
	NewSignature([]string{"arg"}, nil, "", nil, nil, ""),
	func(vs ...Value) Value { return vs[0] })

// TypeOf returns a type name of an argument as a string.
var TypeOf = NewLazyFunction(
	NewSignature([]string{"arg"}, nil, "", nil, nil, ""),
	typeOf)

func typeOf(vs ...Value) Value {
	// No case of effectType should be here.
	switch v := EvalPure(vs[0]).(type) {
	case *BoolType:
		return NewString("bool")
	case *DictionaryType:
		return NewString("dict")
	case *ListType:
		return NewString("list")
	case NilType:
		return NewString("nil")
	case *NumberType:
		return NewString("number")
	case StringType:
		return NewString("string")
	case FunctionType:
		return NewString("function")
	case ErrorType:
		return v
	}

	panic("Unreachable")
}
