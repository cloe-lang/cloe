package core

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestError(t *testing.T) {
	err := PApp(Error, NewString("MyError"), NewString("This is my message.")).Eval().(ErrorType)
	assert.Equal(t, "MyError", err.name)
}

func TestErrorInvalidName(t *testing.T) {
	err := PApp(Error, Nil, NewString("This is my message.")).Eval().(ErrorType)
	assert.Equal(t, "TypeError", err.name)
}

func TestErrorInvalidMessage(t *testing.T) {
	err := PApp(Error, NewString("MyError"), Nil).Eval().(ErrorType)
	assert.Equal(t, "TypeError", err.name)
}
