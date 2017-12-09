package compile

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewModulesCache(t *testing.T) {
	newModulesCache()
}

func TestModulesSet(t *testing.T) {
	assert.Nil(t, newModulesCache().Set("foo", nil))
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
