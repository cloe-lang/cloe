package std

import (
	"github.com/tisp-lang/tisp/src/lib/core"
)

// Map applies a function to each element in a list.
var Map = core.PApp(Y, core.NewLazyFunction(
	core.NewSignature(
		[]string{"self", "func", "list"}, nil, "",
		nil, nil, "",
	),
	func(ts ...*core.Thunk) core.Value {
		f := ts[1]
		l := ts[2]
		return core.PApp(core.If,
			core.PApp(core.Equal, l, core.EmptyList),
			core.EmptyList,
			core.PApp(core.Prepend,
				core.PApp(f, core.PApp(core.First, l)),
				core.PApp(ts[0], f, core.PApp(core.Rest, l))))
	}))
