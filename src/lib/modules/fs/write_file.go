package fs

import (
	"io/ioutil"
	"os"

	"github.com/cloe-lang/cloe/src/lib/core"
)

var writeFile = core.NewEffectFunction(
	core.NewSignature(
		[]string{"file", "content"}, "",
		[]core.OptionalParameter{core.NewOptionalParameter("mode", core.NewNumber(0600))}, "",
	),
	func(vs ...core.Value) core.Value {
		f, err := core.EvalString(vs[0])

		if err != nil {
			return err
		}

		s, err := core.EvalString(vs[1])

		if err != nil {
			return err
		}

		m, err := core.EvalNumber(vs[2])

		if err != nil {
			return err
		}

		if err := ioutil.WriteFile(string(f), []byte(s), os.FileMode(m)); err != nil {
			return fileSystemError(err)
		}

		return core.Nil
	})
