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

func TestPartialError(t *testing.T) {
	for _, a := range []Arguments{
		NewPositionalArguments(Nil),
		NewPositionalArguments(Prepend),
	} {
		_, ok := PApp(App(Partial, a)).Eval().(ErrorType)
		assert.True(t, ok)
	}
}

func TestClosureToString(t *testing.T) {
	assert.Equal(t, StringType("<function>"), PApp(ToString, PApp(Partial, If, True)).Eval())
}
