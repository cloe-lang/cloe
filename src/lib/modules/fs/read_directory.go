package fs

import (
	"io/ioutil"

	"github.com/cloe-lang/cloe/src/lib/core"
)

var readDirectory = core.NewEffectFunction(
	core.NewSignature([]string{"name"}, "", nil, ""),
	func(vs ...core.Value) core.Value {
		s, e := core.EvalString(vs[0])

		if e != nil {
			return e
		}

		fs, err := ioutil.ReadDir(string(s))

		if err != nil {
			return fileSystemError(err)
		}

		ss := make([]core.Value, 0, len(fs))

		for _, f := range fs {
			ss = append(ss, core.NewString(f.Name()))
		}

		return core.NewList(ss...)
	})
