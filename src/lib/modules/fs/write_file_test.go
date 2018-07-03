package fs

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/cloe-lang/cloe/src/lib/core"
	"github.com/stretchr/testify/assert"
)

func TestWriteFile(t *testing.T) {
	d, err := ioutil.TempDir("", "")
	assert.Nil(t, err)

	f := filepath.Join(d, "foo")

	assert.Equal(t,
		core.Nil,
		core.EvalImpure(core.PApp(writeFile, core.NewString(f), core.NewString("bar"))))

	s, err := ioutil.ReadFile(f)
	assert.Nil(t, err)
	assert.Equal(t, "bar", string(s))

	i, err := os.Stat(f)
	assert.Nil(t, err)
	assert.Equal(t, os.FileMode(0600), i.Mode())

	os.RemoveAll(d)
}
