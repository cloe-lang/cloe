package fs

import (
	"io/ioutil"
	"testing"

	"github.com/cloe-lang/cloe/src/lib/core"
	"github.com/stretchr/testify/assert"
)

func TestRemove(t *testing.T) {
	f, err := ioutil.TempFile("", "")
	assert.Nil(t, err)

	_, ok := core.EvalImpure(core.PApp(remove, core.NewString(f.Name()))).(core.NilType)
	assert.True(t, ok)

	e, ok := core.EvalImpure(core.PApp(remove, core.NewString(f.Name()))).(*core.ErrorType)
	assert.True(t, ok)
	assert.Equal(t, "FileSystemError", e.Name())
}
