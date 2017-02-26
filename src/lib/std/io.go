package std

import (
	"github.com/raviqqe/tisp/src/lib/core"
	"fmt"
)

var Write = core.NewStrictFunction(
	core.NewSignature(
		[]string{"x"}, []core.OptionalArgument{}, "",
		[]string{}, []core.OptionalArgument{}, "",
	),
	func(os ...core.Object) core.Object {
		fmt.Println(os[0])

		return core.Nil
	})
