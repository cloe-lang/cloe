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
		v := ts[0].Eval()

		switch x := v.(type) {
		case ErrorType:
			return x
		case dumpable:
			v = x.dump()
		default:
			v = PApp(ToString, Normal(v)).Eval()
		}

		if _, ok := v.(StringType); !ok {
			return NotStringError(v)
		}

		return v
	})

// internalDumpOrFail is the same as DumpOrFail.
func internalDumpOrFail(v Value) string {
	v = ensureWHNF(v)

	switch x := v.(type) {
	case dumpable:
		v = x.dump()
	case stringable:
		v = x.string()
	default:
		panic(fmt.Sprintf("Invalid value detected: %#v", v))
	}

	if s, ok := v.(StringType); ok {
		return string(s)
	}

	panic(fmt.Sprintf("Invalid value detected: %#v", v))
}

// DumpOrFail dumps a value into a string value or fail exiting a process.
// This function should be used only to create strings of error information.
func DumpOrFail(v Value) string {
	v = PApp(Dump, Normal(v)).Eval()
	s, ok := v.(StringType)

	if !ok {
		panic(NotStringError(v).Eval().(ErrorType))
	}

	return string(s)
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

// ensureWHNF evaluates nested thunks into WHNF values.
// This function must be used with care because it prevents tail call
// elimination.
func ensureWHNF(v Value) Value {
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
		// No case of OutputType should be here.
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
		case closureType:
			return NewString("function")

		case ErrorType:
			return NewString("error")
		}

		panic(fmt.Errorf("Invalid value: %#v", ts[0].Eval()))
	})
