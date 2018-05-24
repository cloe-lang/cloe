package random

import (
	"math/rand"

	"github.com/cloe-lang/cloe/src/lib/core"
)

var number = core.NewLazyFunction(
	core.NewSignature(nil, "", nil, ""),
	func(vs ...core.Value) core.Value {
		return core.NewNumber(rand.Float64())
	})
