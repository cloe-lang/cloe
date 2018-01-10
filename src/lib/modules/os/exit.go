package os

import (
	"os"

	"github.com/coel-lang/coel/src/lib/core"
)

var exit = createExitFunction(os.Exit)

func createExitFunction(exit func(int)) core.Value {
	return core.NewEffectFunction(
		core.NewSignature(
			nil, []core.OptionalArgument{core.NewOptionalArgument("status", core.NewNumber(0))}, "",
			nil, nil, ""),
		func(ts ...core.Value) core.Value {
			n, err := ts[0].EvalNumber()

			if err != nil {
				return err
			}

			exit(int(n))

			return core.Nil
		})
}
