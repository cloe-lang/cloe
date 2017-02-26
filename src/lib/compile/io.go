package compile

import (
	"../core"
	"fmt"
)

var write = core.NewStrictFunction(
	core.NewSignature(
		[]string{"x"}, []core.OptionalArgument{}, "",
		[]string{}, []core.OptionalArgument{}, "",
	),
	func(os ...core.Object) core.Object {
		fmt.Println(os[0])

		return core.Nil
	})
