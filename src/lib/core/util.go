package core

import (
	"fmt"

	"github.com/tisp-lang/tisp/src/lib/util"
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

// internalDumpOrFail is the same as DumpOrFail
func internalDumpOrFail(v Value) string {
	v = ensureWHNF(v)

	switch x := v.(type) {
	case ErrorType:
		util.PanicError(x)
	case dumpable:
		v = x.dump()
	case stringable:
		v = x.string()
	default:
		panic(fmt.Sprintf("Invalid value detected: %#v", v))
	}

	switch x := v.(type) {
	case StringType:
		return string(x)
	case ErrorType:
		util.PanicError(x)
	}

	panic(fmt.Sprintf("Invalid value detected: %#v", v))
}

// DumpOrFail dumps a value into a string value or fail exiting a process.
// This function should be used only to create strings of error information.
func DumpOrFail(v Value) string {
	v = PApp(Dump, Normal(v)).Eval()
	s, ok := v.(StringType)

	if !ok {
		util.PanicError(NotStringError(v).Eval().(ErrorType))
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
