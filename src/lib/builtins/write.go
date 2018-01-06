package builtins

import (
	"fmt"
	"os"
	"strings"

	"github.com/coel-lang/coel/src/lib/core"
)

// Write writes string representation of arguments to stdout.
var Write = core.NewEffectFunction(
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
		sep, err := evalGoString(ts[1])

		if err != nil {
			return err
		}

		f, err := evalFileArguments(ts[3], ts[4])

		if err != nil {
			return err
		}

		t := ts[0]
		ss := []string{}

		for {
			if b, err := core.IsEmptyList(t); err != nil {
				return err
			} else if b {
				break
			}

			s, err := evalGoString(core.PApp(core.ToString, core.PApp(core.First, t)))

			if err != nil {
				return err
			}

			ss = append(ss, s)
			t = core.PApp(core.Rest, t)
		}

		end, err := evalGoString(ts[2])

		if err != nil {
			return err
		}

		if _, err := fmt.Fprint(f, strings.Join(ss, sep)+end); err != nil {
			return fileError(err)
		}

		return core.Nil
	})

func evalGoString(t *core.Thunk) (string, core.Value) {
	s, err := core.EvalString(t)

	if err != nil {
		return "", err
	}

	return string(s), nil
}

func evalFileArguments(f, m *core.Thunk) (*os.File, core.Value) {
	switch x := f.Eval().(type) {
	case core.StringType:
		m, e := core.EvalNumber(m)

		if e != nil {
			return nil, e
		}

		f, err := os.OpenFile(
			string(x),
			os.O_CREATE|os.O_TRUNC|os.O_WRONLY,
			os.FileMode(m))

		if err != nil {
			return nil, fileError(err)
		}

		return f, nil
	case core.NumberType:
		switch x {
		case 1:
			return os.Stdout, nil
		case 2:
			return os.Stderr, nil
		}
	}

	return nil, core.ValueError(
		"file optional argument's value must be 1 or 2, or a string filename.")
}
