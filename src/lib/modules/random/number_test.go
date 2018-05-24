package random

import (
	"testing"

	"github.com/cloe-lang/cloe/src/lib/core"
	"github.com/stretchr/testify/assert"
)

func TestNumber(t *testing.T) {
	for i := 0; i < 1000; i++ {
		n, err := core.EvalNumber(core.PApp(number))
		assert.True(t, 0 <= n && n < 1)
		assert.Nil(t, err)
	}
}
