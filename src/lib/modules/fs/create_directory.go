package fs

import (
	"os"

	"github.com/coel-lang/coel/src/lib/core"
)

var createDirectory = core.NewLazyFunction(
	core.NewSignature(
		[]string{"name"}, nil, "", nil,
		[]core.OptionalArgument{core.NewOptionalArgument("existOk", core.False)}, "",
	),
	func(ts ...*core.Thunk) core.Value {
		v := ts[0].Eval()
		s, ok := v.(core.StringType)

		if !ok {
			return core.NotStringError(v)
		}

		v = ts[1].Eval()
		b, ok := v.(core.BoolType)

		if !ok {
			return core.NotBoolError(v)
		}

		if b {
			if f, err := os.Stat(string(s)); err == nil && f.IsDir() {
				return core.NewEffect(core.Nil)
			}
		}

		err := os.Mkdir(string(s), os.ModeDir|0775)

		if err != nil {
			return fileSystemError(err)
		}

		return core.NewEffect(core.Nil)
	})
