package fs

import (
	"io/ioutil"
	"os"
	"path"
	"testing"

	"github.com/cloe-lang/cloe/src/lib/core"
	"github.com/stretchr/testify/assert"
)

func TestReadDirectory(t *testing.T) {
	d, err := ioutil.TempDir("", "")
	assert.Nil(t, err)

	l, e := core.EvalList(core.EvalImpure(core.PApp(readDirectory, core.NewString(d))))
	assert.Nil(t, e)

	n, e := core.EvalNumber(core.PApp(core.Size, l))
	assert.Nil(t, e)
	assert.Equal(t, core.NumberType(0), n)

	f, err := os.OpenFile(path.Join(d, "foo.txt"), os.O_CREATE, 0600)
	assert.Nil(t, err)

	l, e = core.EvalList(core.EvalImpure(core.PApp(readDirectory, core.NewString(d))))
	assert.Nil(t, e)

	n, e = core.EvalNumber(core.PApp(core.Size, l))
	assert.Nil(t, e)
	assert.Equal(t, core.NumberType(1), n)

	assert.Nil(t, os.Remove(f.Name()))
	assert.Nil(t, os.Remove(d))
}

func TestReadDirectoryError(t *testing.T) {
	for _, v := range []core.Value{core.Nil, core.NewString("dir")} {
		_, e := core.EvalList(core.EvalImpure(core.PApp(readDirectory, v)))
		assert.NotNil(t, e)
	}
}
