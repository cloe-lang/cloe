package fs

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/cloe-lang/cloe/src/lib/core"
	"github.com/stretchr/testify/assert"
)

func TestCreateDirectory(t *testing.T) {
	root, err := ioutil.TempDir("", "")

	assert.Nil(t, err)

	d := filepath.Join(root, "foo")

	_, ok := core.EvalImpure(core.PApp(createDirectory, core.NewString(d))).(core.NilType)
	assert.True(t, ok)

	e, ok := core.EvalImpure(core.PApp(createDirectory, core.NewString(d))).(*core.ErrorType)
	assert.True(t, ok)
	assert.Equal(t, "FileSystemError", e.Name())

	_, ok = core.EvalImpure(core.App(
		createDirectory,
		core.NewArguments(
			[]core.PositionalArgument{core.NewPositionalArgument(core.NewString(d), false)},
			[]core.KeywordArgument{core.NewKeywordArgument("existOk", core.True)}),
	)).(core.NilType)
	assert.True(t, ok)

	os.Remove(root)
}

func TestCreateDirectoryWithInvalidArguments(t *testing.T) {
	_, ok := core.EvalImpure(core.PApp(createDirectory, core.Nil)).(*core.ErrorType)
	assert.True(t, ok)

	_, ok = core.EvalImpure(core.App(
		createDirectory,
		core.NewArguments(
			[]core.PositionalArgument{core.NewPositionalArgument(core.NewString("foo"), false)},
			[]core.KeywordArgument{core.NewKeywordArgument("existOk", core.Nil)}),
	)).(*core.ErrorType)
	assert.True(t, ok)

}
