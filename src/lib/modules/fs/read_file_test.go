package fs

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/cloe-lang/cloe/src/lib/core"
	"github.com/stretchr/testify/assert"
)

func TestReadFile(t *testing.T) {
	f, err := ioutil.TempFile("", "")
	assert.Nil(t, err)

	f.Write([]byte("hello"))

	assert.Equal(t,
		core.NewString("hello"),
		core.EvalPure(core.PApp(readFile, core.NewString(f.Name()))))

	os.Remove(f.Name())
}

func TestReadFileError(t *testing.T) {
	for _, v := range []core.Value{
		core.True,
		core.NewString("nonExistentFile"),
		core.DummyError,
	} {
		_, ok := core.EvalPure(core.PApp(readFile, v)).(*core.ErrorType)
		assert.True(t, ok)
	}
}
