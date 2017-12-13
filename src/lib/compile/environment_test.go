package compile

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/coel-lang/coel/src/lib/core"
)

func testFallback(s string) (*core.Thunk, error) {
	return nil, errors.New("name not found")
}

func TestNewEnvironment(t *testing.T) {
	newEnvironment(testFallback)
}

func TestEnvironmentGetFail(t *testing.T) {
	defer func() {
		_, ok := recover().(error)
		assert.True(t, ok)
	}()

	newEnvironment(testFallback).get("foo")
}

func TestEnvironmentGet(t *testing.T) {
	e := newEnvironment(testFallback)
	e.set("foo", core.Nil)
	th := e.get("foo")
	assert.Equal(t, core.Nil, th)
}
