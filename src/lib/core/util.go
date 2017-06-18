package core

import "fmt"

func sprint(x interface{}) StringType {
	return StringType(fmt.Sprint(x))
}

type dumpable interface {
	dump() Value
}

func dump(v Value) (string, Value) {
	if err, ok := v.(ErrorType); ok {
		return "", err
	}

	if d, ok := v.(dumpable); ok {
		v = d.dump()
	} else {
		v = PApp(ToString, Normal(v)).Eval()
	}

	s, ok := v.(StringType)
	if !ok {
		return "", NotStringError(v)
	}

	return string(s), nil
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
