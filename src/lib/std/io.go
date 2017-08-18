package std

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/tisp-lang/tisp/src/lib/core"
)

// Read reads a string from stdin or a file.
var Read = core.NewLazyFunction(
	core.NewSignature(
		nil, []core.OptionalArgument{core.NewOptionalArgument("file", core.Nil)}, "",
		nil, nil, "",
	),
	func(ts ...*core.Thunk) core.Value {
		v := ts[0].Eval()
		file := os.Stdin

		if s, ok := v.(core.StringType); ok {
			var err error
			file, err = os.Open(string(s))

			if err != nil {
				return core.OutputError(err.Error())
			}
		} else if _, ok := v.(core.NilType); !ok {
			return core.ValueError(
				"file optional argument's value must be nil or a filename. Got %s.",
				core.DumpOrFail(v))
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
	func(ts ...*core.Thunk) core.Value {
		v := ts[0].Eval()
		l, ok := v.(core.ListType)

		if !ok {
			return core.NotListError(v)
		}

		elems, err := l.ToValues()

		if err != nil {
			return err
		}

		ss := make([]string, 0, len(elems))

		for _, t := range elems {
			v := core.PApp(core.ToString, t).Eval()
			s, ok := v.(core.StringType)

			if !ok {
				return core.NotStringError(v)
			}

			ss = append(ss, string(s))
		}

		var options [2]string

		for i, t := range ts[1:3] {
			v := t.Eval()
			s, ok := v.(core.StringType)

			if !ok {
				return core.NotStringError(v)
			}

			options[i] = string(s)
		}

		file := os.Stdout

		fileArg := ts[3].Eval()
		if s, ok := fileArg.(core.StringType); ok {
			v := ts[4].Eval()
			mode, ok := v.(core.NumberType)
			if !ok {
				return core.NotNumberError(v)
			}

			var err error
			file, err = os.OpenFile(
				string(s),
				os.O_CREATE|os.O_TRUNC|os.O_WRONLY,
				os.FileMode(mode))

			if err != nil {
				return core.OutputError(err.Error())
			}
		} else if n, ok := fileArg.(core.NumberType); ok && n == 2 {
			file = os.Stderr
		} else if !(ok && n == 1) {
			return core.ValueError(
				"file optional argument's value must be 1 or 2, or a string filename. Got %s.",
				core.DumpOrFail(fileArg))
		}

		fmt.Fprint(file, strings.Join(ss, options[0])+options[1])

		return core.NewOutput(core.Nil)
	})
