package ir

import (
	"testing"

	"github.com/cloe-lang/cloe/src/lib/core"
	"github.com/cloe-lang/cloe/src/lib/debug"
)

func TestAppInterpret(t *testing.T) {
	NewApp(
		core.ToString,
		NewArguments(
			[]PositionalArgument{NewPositionalArgument(core.Nil, false)},
			[]KeywordArgument{NewKeywordArgument("foo", core.Nil), NewKeywordArgument("", core.Nil)}),
		debug.NewGoInfo(0)).interpret(nil)
}
