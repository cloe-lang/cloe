package std

import (
	"testing"

	"github.com/raviqqe/tisp/src/lib/core"
	"github.com/stretchr/testify/assert"
)

func TestSeq(t *testing.T) {
	assert.Equal(
		t,
		core.Nil.Eval(),
		core.PApp(
			Seq,
			core.PApp(Write, core.NewNumber(42)),
			core.PApp(Write, core.NewString("OK!"))).EvalOutput())
}
