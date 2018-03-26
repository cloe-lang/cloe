package ir

import (
	"testing"

	"github.com/coel-lang/coel/src/lib/core"
	"github.com/coel-lang/coel/src/lib/debug"
)

func TestAppInterpret(t *testing.T) {
	NewApp(
		core.ToString,
		NewArguments(
			[]PositionalArgument{NewPositionalArgument(core.Nil, false)},
			[]KeywordArgument{NewKeywordArgument("foo", core.Nil), NewKeywordArgument("", core.Nil)}),
		debug.NewGoInfo(0)).interpret(nil)
}
