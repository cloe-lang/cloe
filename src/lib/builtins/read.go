package builtins

import (
	"io"
	"io/ioutil"
	"os"

	"github.com/coel-lang/coel/src/lib/core"
)

// Read reads a string from stdin or a file.
var Read = createReadFunction(os.Stdin)

func createReadFunction(stdin io.Reader) core.Value {
	return core.NewLazyFunction(
		core.NewSignature(
			nil, []core.OptionalArgument{core.NewOptionalArgument("file", core.Nil)}, "",
			nil, nil, "",
		),
		func(ts ...core.Value) core.Value {
			file := stdin

			switch x := ts[0].Eval().(type) {
			case core.StringType:
				var err error
				file, err = os.Open(string(x))

				if err != nil {
					return fileError(err)
				}
			case core.NilType:
			default:
				s, err := core.StrictDump(x)

				if err != nil {
					return err
				}

				return core.ValueError(
					"file optional argument's value must be nil or a filename. Got %s.",
					s)
			}

			s, err := ioutil.ReadAll(file)

			if err != nil {
				return fileError(err)
			}

			return core.NewString(string(s))
		})
}
