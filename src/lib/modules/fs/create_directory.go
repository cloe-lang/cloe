package fs

import (
	"os"

	"github.com/coel-lang/coel/src/lib/core"
)

var createDirectory = core.NewEffectFunction(
	core.NewSignature(
		[]string{"name"}, nil, "", nil,
		[]core.OptionalArgument{core.NewOptionalArgument("existOk", core.False)}, "",
	),
	func(ts ...core.Value) core.Value {
		s, e := ts[0].EvalString()

		if e != nil {
			return e
		}

		b, e := ts[1].EvalBool()

		if e != nil {
			return e
		}

		if b {
			if f, err := os.Stat(string(s)); err == nil && f.IsDir() {
				return core.Nil
			}
		}

		if err := os.Mkdir(string(s), os.ModeDir|0775); err != nil {
			return fileSystemError(err)
		}

		return core.Nil
	})
