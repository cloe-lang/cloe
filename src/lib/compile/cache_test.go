package compile

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewModulesCache(t *testing.T) {
	assert.NotEqual(t, module(nil), newModulesCache())
}

func TestModulesSet(t *testing.T) {
	assert.Nil(t, newModulesCache().Set("foo", nil))
}

func TestModulesSetInNonExsitentDirectory(t *testing.T) {
	inNonExistentDirectory(t, func() {
		err := newModulesCache().Set("foo", nil)
		assert.NotNil(t, err)
	})
}

func TestModulesGet(t *testing.T) {
	c := newModulesCache()

	err := c.Set("foo", nil)

	assert.Nil(t, err)

	m, ok, err := c.Get("foo")

	assert.Equal(t, module(nil), m)
	assert.True(t, ok)
	assert.Nil(t, err)
}

func TestModulesGetInNonExsitentDirectory(t *testing.T) {
	inNonExistentDirectory(t, func() {
		_, _, err := newModulesCache().Get("foo")
		assert.NotNil(t, err)
	})
}

func inNonExistentDirectory(t *testing.T, f func()) {
	d, err := ioutil.TempDir("", "")
	assert.Nil(t, err)

	wd, err := os.Getwd()
	assert.Nil(t, err)

	os.Chdir(d)
	os.Remove(d)

	f()

	os.Chdir(wd)
}
