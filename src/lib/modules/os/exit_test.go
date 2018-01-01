package os

import (
	"testing"

	"github.com/coel-lang/coel/src/lib/core"
	"github.com/stretchr/testify/assert"
)

func TestExit(t *testing.T) {
	s := 0
	f := createExitFunction(func(i int) { s = i })

	for i, th := range []*core.Thunk{
		core.PApp(f),
		core.PApp(f, core.NewNumber(1)),
		core.PApp(f, core.NewNumber(2)),
	} {
		assert.Equal(t, core.Nil.Eval(), th.EvalEffect())
		assert.Equal(t, i, s)
	}
}

func TestExitError(t *testing.T) {
	f := createExitFunction(func(int) {})

	for _, v := range []core.Value{
		core.PApp(f).Eval(),
		core.PApp(f, core.Nil).EvalEffect(),
	} {
		_, ok := v.(core.ErrorType)
		assert.True(t, ok)
	}
}
