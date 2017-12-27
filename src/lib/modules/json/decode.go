package json

import (
	"github.com/Jeffail/gabs"
	"github.com/coel-lang/coel/src/lib/core"
)

var decode = core.NewLazyFunction(
	core.NewSignature([]string{"encoded"}, nil, "", nil, nil, ""),
	func(ts ...*core.Thunk) core.Value {
		v := ts[0].Eval()
		s, ok := v.(core.StringType)

		if !ok {
			return core.NotStringError(v)
		}

		return decodeString(string(s))
	})

func decodeString(s string) *core.Thunk {
	j, err := gabs.ParseJSON([]byte(s))

	if err != nil {
		return jsonError(err)
	}

	return convertToValue(j.Data())
}

func convertToValue(x interface{}) *core.Thunk {
	if x == nil {
		return core.Nil
	}

	switch x := x.(type) {
	case []interface{}:
		ts := []*core.Thunk{}

		for _, y := range x {
			ts = append(ts, convertToValue(y))
		}

		return core.NewList(ts...)
	case map[string]interface{}:
		kvs := make([]core.KeyValue, 0, len(x))

		for k, v := range x {
			kvs = append(kvs, core.KeyValue{core.NewString(k), convertToValue(v)})
		}

		return core.NewDictionary(kvs)
	case string:
		return core.NewString(x)
	case float64:
		return core.NewNumber(x)
	case bool:
		return core.NewBool(x)
	}

	panic("Unreachable")
}
