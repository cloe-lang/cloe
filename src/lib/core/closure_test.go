package core

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPartial(t *testing.T) {
	ifFunc := func(ts ...*Thunk) bool {
		b := PApp(PApp(Partial, If, False, True), ts...)
		return bool(b.Eval().(BoolType))
	}

	assert.True(t, ifFunc(True))
	assert.True(t, !ifFunc(False))
}

func TestClosureToString(t *testing.T) {
	assert.Equal(t, NewString("<function>").Eval(), PApp(ToString, PApp(Partial, If, True)).Eval())
}
