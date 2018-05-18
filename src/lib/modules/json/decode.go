package json

import (
	"github.com/Jeffail/gabs"
	"github.com/cloe-lang/cloe/src/lib/core"
)

var decode = core.NewLazyFunction(
	core.NewSignature([]string{"encoded"}, "", nil, ""),
	func(vs ...core.Value) core.Value {
		s, err := core.EvalString(vs[0])

		if err != nil {
			return err
		}

		return decodeString(string(s))
	})

func decodeString(s string) core.Value {
	j, err := gabs.ParseJSON([]byte(s))

	if err != nil {
		return jsonError(err)
	}

	return convertToValue(j.Data())
}

func convertToValue(x interface{}) core.Value {
	if x == nil {
		return core.Nil
	}

	switch x := x.(type) {
	case []interface{}:
		ts := []core.Value{}

		for _, y := range x {
			ts = append(ts, convertToValue(y))
		}

		return core.NewList(ts...)
	case map[string]interface{}:
		kvs := make([]core.KeyValue, 0, len(x))

		for k, v := range x {
			kvs = append(kvs, core.KeyValue{Key: core.NewString(k), Value: convertToValue(v)})
		}

		return core.NewDictionary(kvs)
	case string:
		return core.NewString(x)
	case float64:
		return core.NewNumber(x)
	case bool:
		return core.NewBoolean(x)
	}

	panic("Unreachable")
}
