package os

import (
	"strings"
	"testing"

	"github.com/cloe-lang/cloe/src/lib/core"
	"github.com/stretchr/testify/assert"
)

func TestReadStdin(t *testing.T) {
	for _, s := range []string{"", "foo"} {
		assert.Equal(t,
			core.NewString(s),
			core.EvalPure(core.PApp(createReadStdin(strings.NewReader(s)))))
	}
}

func TestReadAsList(t *testing.T) {
	for _, s := range []string{"", "foo"} {
		l, err := core.EvalList(core.App(
			createReadStdin(strings.NewReader(s)),
			core.NewArguments(nil, []core.KeywordArgument{core.NewKeywordArgument("list", core.True)})))
		assert.Nil(t, err)

		for _, r := range s {
			s, err := core.EvalString(l.First())
			assert.Equal(t, core.NewString(string(r)), s)
			assert.Nil(t, err)

			l, err = core.EvalList(l.Rest())
			assert.Nil(t, err)
		}

		assert.Equal(t, core.EmptyList, l)
	}
}
