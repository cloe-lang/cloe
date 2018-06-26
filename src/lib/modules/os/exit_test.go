package os

import (
	"testing"

	"github.com/cloe-lang/cloe/src/lib/core"
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
	_, ok := core.EvalImpure(core.App(
		createExitFunction(func(int) {}),
		core.NewArguments(
			nil,
			[]core.KeywordArgument{core.NewKeywordArgument("status", core.Nil)}))).(*core.ErrorType)

	assert.True(t, ok)
}
