package std

import (
	"github.com/raviqqe/tisp/src/lib/core"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCause(t *testing.T) {
	assert.Equal(
		t,
		core.Nil.Eval(),
		core.PApp(
			Cause,
			core.PApp(Write, core.NewNumber(42)),
			core.PApp(Write, core.NewString("OK!"))).Eval())
}
