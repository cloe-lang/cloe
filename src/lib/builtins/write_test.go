package builtins

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/coel-lang/coel/src/lib/core"
	"github.com/stretchr/testify/assert"
)

func TestWrite(t *testing.T) {
	for _, a := range []core.Arguments{
		core.NewPositionalArguments(core.Nil),
		core.NewPositionalArguments(core.Nil, core.Nil),
		core.NewArguments(
			[]core.PositionalArgument{core.NewPositionalArgument(core.Nil, false)},
			[]core.KeywordArgument{core.NewKeywordArgument("sep", core.NewString(","))},
			nil),
		core.NewArguments(
			[]core.PositionalArgument{core.NewPositionalArgument(core.Nil, false)},
			[]core.KeywordArgument{core.NewKeywordArgument("end", core.NewString(""))},
			nil),
		core.NewArguments(
			[]core.PositionalArgument{core.NewPositionalArgument(core.Nil, false)},
			[]core.KeywordArgument{core.NewKeywordArgument("file", core.NewString("/tmp/coel.test"))},
			nil),
		core.NewArguments(
			[]core.PositionalArgument{core.NewPositionalArgument(core.Nil, false)},
			[]core.KeywordArgument{core.NewKeywordArgument("file", core.NewNumber(2))},
			nil),
		core.NewArguments(
			[]core.PositionalArgument{core.NewPositionalArgument(core.Nil, false)},
			[]core.KeywordArgument{core.NewKeywordArgument("mode", core.NewNumber(0775))},
			nil),
	} {
		assert.Equal(t, core.Nil, core.EvalImpure(core.App(Write, a)))
	}
}

func TestWriteError(t *testing.T) {
	d, err := ioutil.TempDir("", "")
	assert.Nil(t, err)
	defer os.Remove(d)

	for _, a := range []core.Arguments{
		core.NewPositionalArguments(core.DummyError),
		core.NewArguments(
			[]core.PositionalArgument{core.NewPositionalArgument(core.DummyError, true)},
			nil,
			nil),
		core.NewArguments(
			[]core.PositionalArgument{
				core.NewPositionalArgument(core.StrictPrepend([]core.Value{core.Nil}, core.DummyError), true),
			},
			nil,
			nil),
		core.NewPositionalArguments(core.PApp(core.Prepend, core.Nil, core.DummyError)),
		core.NewArguments(
			[]core.PositionalArgument{core.NewPositionalArgument(core.Nil, false)},
			[]core.KeywordArgument{core.NewKeywordArgument("sep", core.Nil)},
			nil),
		core.NewArguments(
			[]core.PositionalArgument{core.NewPositionalArgument(core.Nil, false)},
			[]core.KeywordArgument{core.NewKeywordArgument("end", core.Nil)},
			nil),
		core.NewArguments(
			[]core.PositionalArgument{core.NewPositionalArgument(core.Nil, false)},
			[]core.KeywordArgument{
				core.NewKeywordArgument("file", core.NewString("newFile")),
				core.NewKeywordArgument("mode", core.Nil),
			},
			nil),
		core.NewArguments(
			[]core.PositionalArgument{core.NewPositionalArgument(core.Nil, false)},
			[]core.KeywordArgument{
				core.NewKeywordArgument("file", core.DummyError),
			},
			nil),
		core.NewArguments(
			[]core.PositionalArgument{core.NewPositionalArgument(core.Nil, false)},
			[]core.KeywordArgument{
				core.NewKeywordArgument("file", core.True),
			},
			nil),
		core.NewArguments(
			[]core.PositionalArgument{core.NewPositionalArgument(core.Nil, false)},
			[]core.KeywordArgument{
				core.NewKeywordArgument("file", core.True),
			},
			nil),
		core.NewArguments(
			[]core.PositionalArgument{core.NewPositionalArgument(core.Nil, false)},
			[]core.KeywordArgument{
				core.NewKeywordArgument("file", core.NewString("/dev/full")),
			},
			nil),
		core.NewArguments(
			[]core.PositionalArgument{core.NewPositionalArgument(core.Nil, false)},
			[]core.KeywordArgument{
				core.NewKeywordArgument("file", core.NewString("/dev/full")),
			},
			nil),
		core.NewArguments(
			[]core.PositionalArgument{core.NewPositionalArgument(core.Nil, false)},
			[]core.KeywordArgument{
				core.NewKeywordArgument("file", core.NewString(d)),
			},
			nil),
	} {
		_, ok := core.EvalPure(core.App(Write, a)).(core.ErrorType)
		assert.True(t, ok)
	}
}
