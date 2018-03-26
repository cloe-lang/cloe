package os

import (
	"testing"

	"github.com/coel-lang/coel/src/lib/core"
	"github.com/stretchr/testify/assert"
)

func TestExit(t *testing.T) {
	s := 0
	f := createExitFunction(func(i int) { s = i })

	for i, v := range []core.Value{
		core.PApp(f),
		core.App(f, core.NewArguments(nil, []core.KeywordArgument{core.NewKeywordArgument("status", core.NewNumber(1))})),
		core.App(f, core.NewArguments(nil, []core.KeywordArgument{core.NewKeywordArgument("status", core.NewNumber(2))})),
	} {
		assert.Equal(t, core.Nil, core.EvalImpure(v))
		assert.Equal(t, i, s)
	}
}

func TestExitError(t *testing.T) {
	f := createExitFunction(func(int) {})

	for _, v := range []core.Value{
		core.EvalPure(core.PApp(f)),
		core.EvalImpure(core.PApp(f, core.Nil)),
	} {
		_, ok := v.(core.ErrorType)
		assert.True(t, ok)
	}
}
