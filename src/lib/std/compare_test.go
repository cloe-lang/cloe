package std

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tisp-lang/tisp/src/lib/core"
)

func TestLessTrue(t *testing.T) {
	for _, ts := range [][]*core.Thunk{
		{},
		{core.NewNumber(42)},
		{core.NewString("foo")},
		{core.NewList(core.NewNumber(42))},
		{core.NewNumber(0), core.NewNumber(0.0000000000001)},
		{core.NewNumber(0), core.NewNumber(0.00000001), core.NewNumber(0.0001)},
		{core.NewString("bar"), core.NewString("baz"), core.NewString("foo")},
		{core.NewList(core.NewNumber(42)), core.NewList(core.NewNumber(2049))},
		{core.NewList(core.NewNumber(42)), core.NewList(core.NewNumber(42), core.NewNumber(42))},
		{core.NewList(core.NewList(core.NewNumber(42))),
			core.NewList(core.NewList(core.NewNumber(42), core.NewNumber(42))),
			core.NewList(core.NewList(core.NewNumber(42), core.NewNumber(42), core.NewNumber(42)))},
	} {
		assert.True(t, bool(core.PApp(Less, ts...).Eval().(core.BoolType)))
	}
}

func TestLessFalse(t *testing.T) {
	for _, ts := range [][]*core.Thunk{
		{core.True, core.False},
		{core.NewNumber(42), core.NewNumber(0)},
		{core.NewNumber(0), core.NewNumber(42), core.NewNumber(42)},
		{core.NewNumber(0), core.NewNumber(0), core.NewNumber(42)},
		{core.NewString("foo"), core.NewString("bar"), core.NewString("baz")},
		{core.NewList(core.NewList(core.NewNumber(42))),
			core.NewList(core.NewList(core.NewNumber(42), core.NewNumber(42))),
			core.NewList(core.NewList(core.NewNumber(42), core.NewNumber(42)))},
		{core.NewList(core.NewList(core.NewNumber(42), core.NewNumber(42))),
			core.NewList(core.NewList(core.NewNumber(42), core.NewNumber(42))),
			core.NewList(core.NewList(core.NewNumber(42), core.NewNumber(2049)))},
	} {
		assert.True(t, !bool(core.PApp(Equal, ts...).Eval().(core.BoolType)))
	}
}
