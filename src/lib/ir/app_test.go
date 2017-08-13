package ir

import (
	"testing"

	"github.com/tisp-lang/tisp/src/lib/core"
	"github.com/tisp-lang/tisp/src/lib/debug"
)

func TestAppInterpret(t *testing.T) {
	NewApp(core.ToString, NewArguments(
		[]PositionalArgument{NewPositionalArgument(core.Nil, false)},
		[]KeywordArgument{NewKeywordArgument("foo", core.Nil)},
		[]interface{}{core.Nil}), debug.NewGoInfo(0)).interpret(nil)
}
