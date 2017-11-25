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
	NewSignature([]string{"x"}, nil, "", nil, nil, ""),
	func(ts ...*Thunk) Value {
		s, err := strictDump(ts[0].Eval())

		if err != nil {
			return err
		}

		return s
	})

func strictDump(v Value) (StringType, *Thunk) {
	switch x := v.(type) {
	case ErrorType:
		return "", Normal(x)
	case dumpable:
		v = x.dump()
	case stringable:
		v = x.string()
	default:
		panic(fmt.Errorf("Invalid value: %#v", x))
	}

	s, ok := v.(StringType)

	if !ok {
		return "", NotStringError(v)
	}

	return s, nil
}

// StrictDump is a variant of Dump which evaluates input strictly.
func StrictDump(v Value) (string, *Thunk) {
	s, err := strictDump(ensureNormal(v))
	return string(s), err
}

// Equal checks if 2 values are equal or not.
var Equal = NewStrictFunction(
	NewSignature(
		[]string{"x", "y"}, nil, "",
		nil, nil, "",
	),
	func(ts ...*Thunk) (v Value) {
		defer func() {
			if r := recover(); r != nil {
				v = r
			}
		}()

		return BoolType(compare(ts[0].Eval(), ts[1].Eval()) == 0)
	})

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
	NewSignature([]string{"x"}, nil, "", nil, nil, ""),
	func(ts ...*Thunk) Value { return ts[0] })

// TypeOf returns a type name of an argument as a string.
var TypeOf = NewLazyFunction(
	NewSignature([]string{"x"}, nil, "", nil, nil, ""),
	func(ts ...*Thunk) Value {
		// No case of EffectType should be here.
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
		case rawFunction:
			return NewString("function")

		case ErrorType:
			return NewString("error")
		}

		panic(fmt.Errorf("Invalid value: %#v", ts[0].Eval()))
	})
