package std

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/raviqqe/tisp/src/lib/core"
)

// Read reads a string from stdin or a file.
var Read = core.NewStrictFunction(
	core.NewSignature(
		nil, []core.OptionalArgument{core.NewOptionalArgument("file", core.Nil)}, "",
		nil, nil, "",
	),
	func(vs ...core.Value) core.Value {
		v := vs[0]
		file := os.Stdin

		if s, ok := v.(core.StringType); ok {
			var err error
			file, err = os.Open(string(s))

			if err != nil {
				return core.OutputError(err.Error())
			}
		} else if _, ok := v.(core.NilType); !ok {
			return core.ValueError("file optional argument's value must be nil or a filename. Got %#v.", v)
		}

		s, err := ioutil.ReadAll(file)

		if err != nil {
			return core.OutputError(err.Error())
		}

		return core.NewString(string(s))
	})

// Write writes string representation of arguments to stdout.
var Write = core.NewStrictFunction(
	core.NewSignature(
		nil, nil, "args",
		nil, []core.OptionalArgument{
			core.NewOptionalArgument("sep", core.NewString(" ")),
			core.NewOptionalArgument("end", core.NewString("\n")),
			core.NewOptionalArgument("file", core.NewNumber(1)),
			core.NewOptionalArgument("mode", core.NewNumber(0664)),
		}, "",
	),
	func(vs ...core.Value) core.Value {
		elems, err := vs[0].(core.ListType).ToValues()

		if err != nil {
			return err
		}

		ss := make([]string, len(elems))

		for i, v := range elems {
			v := core.PApp(core.ToString, core.Normal(v)).Eval()
			s, ok := v.(core.StringType)

			if !ok {
				return core.NotStringError(v)
			}

			ss[i] = string(s)
		}

		var options [2]string

		for i, v := range vs[1:3] {
			s, ok := v.(core.StringType)

			if !ok {
				return core.NotStringError(v)
			}

			options[i] = string(s)
		}

		file := os.Stdout

		if s, ok := vs[3].(core.StringType); ok {
			mode, ok := vs[4].(core.NumberType)

			if !ok {
				return core.NotNumberError(vs[4])
			}

			var err error
			file, err = os.OpenFile(
				string(s),
				os.O_CREATE|os.O_TRUNC|os.O_WRONLY,
				os.FileMode(mode))

			if err != nil {
				return core.OutputError(err.Error())
			}
		} else if n, ok := vs[3].(core.NumberType); ok && n == 2 {
			file = os.Stderr
		} else if !(ok && n == 1) {
			return core.ValueError("file optional argument's value must be 1 or 2, or a string filename. Got %#v.", vs[3])
		}

		fmt.Fprint(file, strings.Join(ss, options[0])+options[1])

		return core.NewOutput(core.Nil)
	})
