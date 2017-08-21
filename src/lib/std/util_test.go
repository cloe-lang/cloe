package std

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tisp-lang/tisp/src/lib/core"
)

func TestCheckEmptyListFail(t *testing.T) {
	v := checkEmptyList(core.OutOfRangeError(), core.Nil)
	t.Log(v)
	_, ok := v.(core.ErrorType)
	assert.True(t, ok)
}
