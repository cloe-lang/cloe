package compile

import (
	"errors"
	"testing"

	"github.com/cloe-lang/cloe/src/lib/core"
	"github.com/stretchr/testify/assert"
)

func testFallback(s string) (core.Value, error) {
	return nil, errors.New("name not found")
}

func TestNewEnvironment(t *testing.T) {
	newEnvironment(testFallback)
}

func TestEnvironmentGet(t *testing.T) {
	e := newEnvironment(testFallback)
	e.set("foo", core.Nil)
	v, err := e.get("foo")
	assert.Nil(t, err)
	assert.Equal(t, core.Nil, core.EvalPure(v))
}

func TestEnvironmentGetError(t *testing.T) {
	_, err := newEnvironment(testFallback).get("foo")
	assert.NotNil(t, err)
}
