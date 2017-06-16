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

var equal = NewStrictFunction(
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
