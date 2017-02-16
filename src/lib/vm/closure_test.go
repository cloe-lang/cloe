package vm

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestPartial(t *testing.T) {
	ifFunc := func(ts ...*Thunk) bool {
		b := PApp(PApp(Partial, If, False, True), ts...)
		return bool(b.Eval().(BoolType))
	}

	assert.True(t, ifFunc(True))
	assert.True(t, !ifFunc(False))
}
