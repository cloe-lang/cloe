package std

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tisp-lang/tisp/src/lib/core"
)

func TestEqualTrue(t *testing.T) {
	for _, ts := range [][]*core.Thunk{
		{},
		{core.True},
		{core.True, core.True},
		{core.Nil, core.Nil},
		{core.NewNumber(42), core.NewNumber(42)},
		{core.NewNumber(42), core.NewNumber(42), core.NewNumber(42)},
	} {
		assert.True(t, bool(core.PApp(Equal, ts...).Eval().(core.BoolType)))
	}
}

func TestEqualFalse(t *testing.T) {
	for _, ts := range [][]*core.Thunk{
		{core.True, core.False},
		{core.NewNumber(0), core.NewNumber(42)},
		{core.NewNumber(0), core.NewNumber(0), core.NewNumber(42)},
		{core.NewNumber(0), core.NewNumber(42), core.NewNumber(42)},
	} {
		assert.True(t, !bool(core.PApp(Equal, ts...).Eval().(core.BoolType)))
	}
}
