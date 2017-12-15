package json

import (
	"github.com/coel-lang/coel/src/lib/core"
)

var encode = core.NewLazyFunction(
	core.NewSignature([]string{"decoded"}, nil, "", nil, nil, ""),
	func(ts ...*core.Thunk) core.Value {
		return ts[0]
	})
