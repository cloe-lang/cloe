package os

import (
	"os"

	"github.com/coel-lang/coel/src/lib/core"
)

var exit = createExitFunction(os.Exit)

func createExitFunction(exit func(int)) *core.Thunk {
	return core.NewEffectFunction(
		core.NewSignature(
			nil, []core.OptionalArgument{core.NewOptionalArgument("status", core.NewNumber(0))}, "",
			nil, nil, ""),
		func(ts ...*core.Thunk) core.Value {
			v := ts[0].Eval()
			n, ok := v.(core.NumberType)

			if !ok {
				return core.NotNumberError(v)
			}

			exit(int(n))

			return core.Nil
		})
}
