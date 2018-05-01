package os

import (
	"os"

	"github.com/cloe-lang/cloe/src/lib/core"
)

var exit = createExitFunction(os.Exit)

func createExitFunction(exit func(int)) core.Value {
	return core.NewEffectFunction(
		core.NewSignature(
			nil, "",
			[]core.OptionalParameter{core.NewOptionalParameter("status", core.NewNumber(0))}, ""),
		func(vs ...core.Value) core.Value {
			n, err := core.EvalNumber(vs[0])

			if err != nil {
				return err
			}

			exit(int(n))

			return core.Nil
		})
}
