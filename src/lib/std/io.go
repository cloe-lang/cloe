package std

import (
	"fmt"
	"github.com/raviqqe/tisp/src/lib/core"
	"strings"
)

// Write writes string representation of arguments to stdout.
var Write = core.NewStrictFunction(
	core.NewSignature(
		[]string{}, []core.OptionalArgument{}, "args",
		[]string{}, []core.OptionalArgument{
			core.NewOptionalArgument("sep", core.NewString(" ")),
			core.NewOptionalArgument("end", core.NewString("\n")),
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

		for i, o := range os[1:] {
			s, ok := o.(core.StringType)

			if !ok {
				return core.NotStringError(o)
			}

			options[i] = string(s)
		}

		fmt.Print(strings.Join(ss, options[0]) + options[1])

		return core.Nil
	})
