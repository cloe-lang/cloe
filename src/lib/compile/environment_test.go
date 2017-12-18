package compile

import (
	"errors"
	"testing"

	"github.com/coel-lang/coel/src/lib/core"
	"github.com/stretchr/testify/assert"
)

func testFallback(s string) (*core.Thunk, error) {
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
	th := e.get("foo")
	assert.Equal(t, core.Nil.Eval(), th.Eval())
}
