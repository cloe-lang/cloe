package fs

import (
	"os"

	"github.com/cloe-lang/cloe/src/lib/core"
)

var createDirectory = core.NewEffectFunction(
	core.NewSignature(
		[]string{"name"}, "",
		[]core.OptionalParameter{core.NewOptionalParameter("existOk", core.False)}, "",
	),
	func(vs ...core.Value) core.Value {
		s, e := core.EvalString(vs[0])

		if e != nil {
			return e
		}

		b, e := core.EvalBool(vs[1])

		if e != nil {
			return e
		} else if b {
			if f, err := os.Stat(string(s)); err == nil && f.IsDir() {
				return core.Nil
			}
		}

		if err := os.Mkdir(string(s), os.ModeDir|0775); err != nil {
			return fileSystemError(err)
		}

		return core.Nil
	})
