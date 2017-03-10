package core

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNilEqual(t *testing.T) {
	assert.True(t, testEqual(Nil, Nil))
	assert.True(t, !testEqual(Nil, True))
}
