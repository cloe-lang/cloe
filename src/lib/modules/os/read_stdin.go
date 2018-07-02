package os

import (
	"bufio"
	"io"
	"io/ioutil"
	"os"

	"github.com/cloe-lang/cloe/src/lib/builtins"
	"github.com/cloe-lang/cloe/src/lib/core"
)

var readStdin = createReadStdin(os.Stdin)

func createReadStdin(r io.Reader) core.Value {
	return core.NewLazyFunction(
		core.NewSignature(
			nil, "",
			[]core.OptionalParameter{core.NewOptionalParameter("list", core.False)}, "",
		),
		func(vs ...core.Value) core.Value {
			b, err := core.EvalBoolean(vs[0])

			if err != nil {
				return err
			}

			if b {
				return readAsList(r)
			}

			return readAsString(r)
		})
}

func readAsString(r io.Reader) core.Value {
	s, err := ioutil.ReadAll(r)

	if err != nil {
		return osError(err)
	}

	return core.NewString(string(s))
}

func readAsList(rd io.Reader) core.Value {
	r := bufio.NewReader(rd)

	f := core.PApp(builtins.Y, core.NewLazyFunction(
		core.NewSignature([]string{"self"}, "", nil, ""),
		func(vs ...core.Value) core.Value {
			r, _, err := r.ReadRune()

			if err != nil && err == io.EOF {
				return core.EmptyList
			} else if err != nil {
				return osError(err)
			}

			return core.StrictPrepend([]core.Value{core.NewString(string(r))}, core.PApp(vs[0]))
		}))

	return core.PApp(f)
}
