package builtins

import "github.com/coel-lang/coel/src/lib/core"

func checkEmptyList(l *core.Thunk, ifTrue core.Value) core.Value {
	b, err := core.EvalBool(core.PApp(core.Equal, l, core.EmptyList))

	if err != nil {
		return err
	} else if b {
		return ifTrue
	}

	return nil
}

func fileError(err error) *core.Thunk {
	return core.NewError("FileError", err.Error())
}
