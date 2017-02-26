package core

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNilEqual(t *testing.T) {
	assert.True(t, testEqual(Nil, Nil))
	assert.True(t, !testEqual(Nil, True))
}
