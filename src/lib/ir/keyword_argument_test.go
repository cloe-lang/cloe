package ir

import (
	"testing"

	"github.com/tisp-lang/tisp/src/lib/core"
)

func TestNewKeywordArgument(t *testing.T) {
	NewKeywordArgument("foo", 0)
}

func TestKeywordArgumentInterpret(t *testing.T) {
	NewKeywordArgument("foo", 0).interpret([]*core.Thunk{core.NewNumber(123)})
}
