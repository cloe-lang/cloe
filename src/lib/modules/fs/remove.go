package fs

import (
	"os"

	"github.com/cloe-lang/cloe/src/lib/core"
)

var remove = core.NewEffectFunction(
	core.NewSignature([]string{"name"}, "", nil, ""),
	func(vs ...core.Value) core.Value {
		s, e := core.EvalString(vs[0])

		if e != nil {
			return e
		}

		if err := os.Remove(string(s)); err != nil {
			return fileSystemError(err)
		}

		return core.Nil
	})
