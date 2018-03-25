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
		nil, "args",
		[]core.OptionalParameter{
			core.NewOptionalParameter("sep", core.NewString(" ")),
			core.NewOptionalParameter("end", core.NewString("\n")),
			core.NewOptionalParameter("file", core.NewNumber(1)),
			core.NewOptionalParameter("mode", core.NewNumber(0664)),
		}, "",
	),
	func(ts ...core.Value) core.Value {
		sep, err := evalGoString(ts[1])

		if err != nil {
			return err
		}

		f, err := evalFileArguments(ts[3], ts[4])

		if err != nil {
			return err
		}

		l, err := core.EvalList(ts[0])

		if err != nil {
			return err
		}

		ss := []string{}

		for !l.Empty() {
			s, err := evalGoString(core.PApp(core.ToString, l.First()))

			if err != nil {
				return err
			}

			ss = append(ss, s)

			l, err = core.EvalList(l.Rest())

			if err != nil {
				return err
			}
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

func evalGoString(t core.Value) (string, core.Value) {
	s, err := core.EvalString(t)

	if err != nil {
		return "", err
	}

	return string(s), nil
}

func evalFileArguments(f, m core.Value) (*os.File, core.Value) {
	switch x := core.EvalPure(f).(type) {
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
	case *core.NumberType:
		switch *x {
		case 1:
			return os.Stdout, nil
		case 2:
			return os.Stderr, nil
		}
	}

	return nil, core.ValueError(
		"file optional argument's value must be 1 or 2, or a string filename.")
}
