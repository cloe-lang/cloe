package core

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestError(t *testing.T) {
	assert.Equal(t,
		"MyError",
		EvalPure(PApp(Error, NewString("MyError"), NewString("This is my message."))).(*ErrorType).Name())
}

func TestErrorInvalidName(t *testing.T) {
	assert.Equal(t,
		"TypeError",
		EvalPure(PApp(Error, Nil, NewString("This is my message."))).(*ErrorType).Name())
}

func TestErrorInvalidMessage(t *testing.T) {
	assert.Equal(t,
		"TypeError",
		EvalPure(PApp(Error, NewString("MyError"), Nil)).(*ErrorType).Name())
}

func TestErrorName(t *testing.T) {
	assert.Equal(t, "DummyError", DummyError.Name())
}

func TestErrorCatch(t *testing.T) {
	_, ok := EvalPure(PApp(Catch, DummyError)).(*DictionaryType)
	assert.True(t, ok)

	_, ok = EvalPure(PApp(Catch, Nil)).(NilType)
	assert.True(t, ok)
}

func TestKeyNotFoundErrorFail(t *testing.T) {
	assert.NotEqual(t,
		"KeyNotFoundError",
		keyNotFoundError(DummyError).Name())
}
