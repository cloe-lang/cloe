package core

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFunctionToString(t *testing.T) {
	assert.Equal(t, StringType("<function>"), PApp(ToString, If).Eval())
}
