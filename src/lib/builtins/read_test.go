package builtins

import (
	"testing"

	"github.com/coel-lang/coel/src/lib/core"
	"github.com/stretchr/testify/assert"
)

func TestReadWithStdin(t *testing.T) {
	assert.Equal(t, core.StringType(""), core.PApp(Read).Eval())
}

func TestReadError(t *testing.T) {
	for _, a := range []core.Arguments{
		core.NewArguments(
			[]core.PositionalArgument{core.NewPositionalArgument(core.True, false)},
			nil,
			nil),
		core.NewArguments(
			[]core.PositionalArgument{
				core.NewPositionalArgument(core.NewString("nonExistentFile"), false),
			},
			nil,
			nil),
		core.NewArguments(
			[]core.PositionalArgument{
				core.NewPositionalArgument(core.NewError("", ""), false),
			},
			nil,
			nil),
	} {
		_, ok := core.App(Read, a).Eval().(core.ErrorType)
		assert.True(t, ok)
	}
}
