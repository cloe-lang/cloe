package builtins

import (
	"testing"

	"github.com/coel-lang/coel/src/lib/core"
	"github.com/stretchr/testify/assert"
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

func TestEqualFail(t *testing.T) {
	for _, th := range []*core.Thunk{
		core.App(
			Equal,
			core.NewArguments(
				[]core.PositionalArgument{core.NewPositionalArgument(core.OutOfRangeError(), true)},
				nil,
				nil)),
		core.App(
			Equal,
			core.NewArguments(
				[]core.PositionalArgument{
					core.NewPositionalArgument(core.Nil, false),
					core.NewPositionalArgument(core.OutOfRangeError(), true),
				},
				nil,
				nil)),
		core.App(
			Equal,
			core.NewArguments(
				[]core.PositionalArgument{
					core.NewPositionalArgument(core.Nil, false),
					core.NewPositionalArgument(core.NewList(core.OutOfRangeError()), true),
				},
				nil,
				nil)),
	} {
		_, ok := th.Eval().(core.ErrorType)
		assert.True(t, ok)
	}
}
