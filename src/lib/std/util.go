package std

import "github.com/tisp-lang/tisp/src/lib/core"

func checkEmptyList(l *core.Thunk, ifTrue core.Value) core.Value {
	v := core.PApp(core.Equal, l, core.EmptyList).Eval()
	b, ok := v.(core.BoolType)

	if !ok {
		return core.NotBoolError(v).Eval()
	} else if b {
		return ifTrue
	}

	return nil
}
