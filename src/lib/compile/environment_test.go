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

func TestEnvironmentGetFail(t *testing.T) {
	defer func() {
		assert.NotNil(t, recover())
	}()

	newEnvironment(testFallback).get("foo")
}

func TestEnvironmentGet(t *testing.T) {
	e := newEnvironment(testFallback)
	e.set("foo", core.Nil)
	v := e.get("foo")
	assert.Equal(t, core.Nil, core.EvalPure(v))
}
