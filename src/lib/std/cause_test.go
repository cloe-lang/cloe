package std

import (
	"../core"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCause(t *testing.T) {
	_, ok := core.PApp(Cause, core.PApp(Write, core.NewNumber(42)), core.PApp(Write, core.NewString("OK!"))).Eval().(core.NilType)
	assert.True(t, ok)
}
