package std

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tisp-lang/tisp/src/lib/core"
)

func TestMap(t *testing.T) {
	th := core.PApp(
		Map,
		core.NewLazyFunction(
			core.NewSignature([]string{"x"}, nil, "", nil, nil, ""),
			func(ts ...*core.Thunk) core.Value {
				return core.PApp(core.Mul, ts[0], ts[0])
			}),
		core.NewList(core.NewNumber(2), core.NewNumber(3)))

	assert.Equal(t, core.NewNumber(4).Eval(), core.PApp(th, core.NewNumber(0)).Eval())
	assert.Equal(t, core.NewNumber(9).Eval(), core.PApp(th, core.NewNumber(1)).Eval())
	assert.Equal(t, core.EmptyList.Eval(), core.PApp(core.Rest, core.PApp(core.Rest, th)).Eval())
}
