package std

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tisp-lang/tisp/src/lib/core"
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
