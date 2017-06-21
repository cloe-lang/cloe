package macro

import (
	"fmt"

	"github.com/tisp-lang/tisp/src/lib/core"
	"github.com/tisp-lang/tisp/src/lib/util"
)

// FunctionsToMacros converts a name-to-function map into a name-to-macro map.
func FunctionsToMacros(fs map[string]*core.Thunk) map[string]func(...interface{}) interface{} {
	gs := map[string]func(...interface{}) interface{}{}

	for k, v := range fs {
		gs[k] = functionToMacro(v)
	}

	return gs
}

func functionToMacro(f *core.Thunk) func(...interface{}) interface{} {
	return func(xs ...interface{}) interface{} {
		ts := make([]*core.Thunk, 0, len(xs))

		for _, x := range xs {
			ts = append(ts, astToThunk(x))
		}

		return thunkToAST(core.PApp(f, ts...))
	}
}

func astToThunk(a interface{}) *core.Thunk {
	switch v := a.(type) {
	case string:
		return core.NewString(v)
	case []interface{}:
		ts := make([]*core.Thunk, 0, len(v))

		for _, a := range v {
			ts = append(ts, astToThunk(a))
		}

		return core.NewList(ts...)
	}

	panic(fmt.Sprintf("invalid type: %#v", a))
}

func thunkToAST(t *core.Thunk) interface{} {
	switch v := t.Eval().(type) {
	case core.StringType:
		return string(v)
	case core.ListType:
		ts, err := v.ToThunks()

		if err != nil {
			util.PanicError(err.Eval().(core.ErrorType))
		}

		as := make([]interface{}, 0, len(ts))

		for _, t := range ts {
			as = append(as, thunkToAST(t))
		}

		return as
	default:
		panic(fmt.Sprintf("invalid type: %#v", v))
	}
}
