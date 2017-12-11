package core

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFunctionToString(t *testing.T) {
	assert.Equal(t, NewString("<function>").Eval(), PApp(ToString, If).Eval())
}
