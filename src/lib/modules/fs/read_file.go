package fs

import (
	"io/ioutil"
	"os"

	"github.com/cloe-lang/cloe/src/lib/core"
)

var readFile = core.NewLazyFunction(
	core.NewSignature([]string{"file"}, "", nil, ""),
	func(vs ...core.Value) core.Value {
		s, e := core.EvalString(vs[0])

		if e != nil {
			return e
		}

		f, err := os.Open(string(s))

		if err != nil {
			return fileSystemError(err)
		}

		bs, err := ioutil.ReadAll(f)

		if err != nil {
			return fileSystemError(err)
		}

		return core.NewString(string(bs))
	})
