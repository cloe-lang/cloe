package core

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPure(t *testing.T) {
	assert.Equal(t, True, EvalPure(PApp(Pure, PApp(impureFunction, True))))
}
