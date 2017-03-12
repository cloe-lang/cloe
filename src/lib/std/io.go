package std

import (
	"fmt"
	osys "os"
	"strings"

	"github.com/raviqqe/tisp/src/lib/core"
)

// Write writes string representation of arguments to stdout.
var Write = core.NewStrictFunction(
	core.NewSignature(
		[]string{}, []core.OptionalArgument{}, "args",
		[]string{}, []core.OptionalArgument{
			core.NewOptionalArgument("sep", core.NewString(" ")),
			core.NewOptionalArgument("end", core.NewString("\n")),
			core.NewOptionalArgument("file", core.NewNumber(1)),
			core.NewOptionalArgument("mode", core.NewNumber(0664)),
		}, "",
	),
	func(os ...core.Object) core.Object {
		elems, err := os[0].(core.ListType).ToObjects()

		if err != nil {
			return err
		}

		ss := make([]string, len(elems))

		for i, o := range elems {
			o := core.PApp(core.ToString, core.Normal(o)).Eval()
			s, ok := o.(core.StringType)

			if !ok {
				return core.NotStringError(o)
			}

			ss[i] = string(s)
		}

		var options [2]string

		for i, o := range os[1:3] {
			s, ok := o.(core.StringType)

			if !ok {
				return core.NotStringError(o)
			}

			options[i] = string(s)
		}

		file := osys.Stdout

		if s, ok := os[3].(core.StringType); ok {
			mode, ok := os[4].(core.NumberType)

			if !ok {
				return core.NotNumberError(os[4])
			}

			var err error
			file, err = osys.OpenFile(
				string(s),
				osys.O_CREATE|osys.O_TRUNC|osys.O_WRONLY,
				osys.FileMode(mode))

			if err != nil {
				return core.OutputError(err.Error())
			}
		} else if n, ok := os[3].(core.NumberType); ok && n == 2 {
			file = osys.Stderr
		} else if !(ok && n == 1) {
			return core.ValueError("file optional argument's value must be 1 or 2, or a string filename. Got %#v.", os[3])
		}

		fmt.Fprint(file, strings.Join(ss, options[0])+options[1])

		return core.Nil
	})
