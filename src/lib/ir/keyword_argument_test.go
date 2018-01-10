package ir

import (
	"testing"

	"github.com/coel-lang/coel/src/lib/core"
)

func TestNewKeywordArgument(t *testing.T) {
	NewKeywordArgument("foo", 0)
}

func TestKeywordArgumentInterpret(t *testing.T) {
	NewKeywordArgument("foo", 0).interpret([]core.Value{core.NewNumber(123)})
}
