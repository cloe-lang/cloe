package builtins

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/coel-lang/coel/src/lib/core"
	"github.com/stretchr/testify/assert"
)

func TestReadWithStdin(t *testing.T) {
	assert.Equal(t, core.StringType(""), core.PApp(Read).Eval())
}

func TestReadError(t *testing.T) {
	for _, th := range []core.Value{
		core.True,
		core.NewString("nonExistentFile"),
		core.DummyError,
	} {
		_, ok := core.PApp(Read, th).Eval().(core.ErrorType)
		assert.True(t, ok)
	}
}

func TestReadWithClosedStdin(t *testing.T) {
	f, err := ioutil.TempFile("", "")
	assert.Nil(t, err)
	err = f.Close()
	assert.Nil(t, err)
	defer os.Remove(f.Name())

	_, ok := core.PApp(createReadFunction(f)).Eval().(core.ErrorType)
	assert.True(t, ok)
}
