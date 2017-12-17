package builtins

import (
	"testing"

	"github.com/coel-lang/coel/src/lib/core"
	"github.com/stretchr/testify/assert"
)

func TestSeq(t *testing.T) {
	for _, ts := range [][]*core.Thunk{
		{core.Nil},
		{core.Nil, core.Nil},
	} {
		assert.Equal(
			t,
			core.Nil.Eval(),
			core.PApp(Seq, ts...).Eval())
	}
}

func TestSeqWithEffects(t *testing.T) {
	w := core.PApp(Write, core.Nil)

	for _, ts := range [][]*core.Thunk{
		{w},
		{w, w},
	} {
		_, ok := core.PApp(Seq, ts...).Eval().(core.ErrorType)
		assert.True(t, ok)
	}
}

func TestEffectSeq(t *testing.T) {
	assert.Equal(
		t,
		core.Nil.Eval(),
		core.PApp(
			EffectSeq,
			core.PApp(Write, core.NewNumber(42)),
			core.PApp(Write, core.NewString("OK!"))).EvalEffect())
}

func TestEffectSeqWithPureValues(t *testing.T) {
	for _, ts := range [][]*core.Thunk{
		{core.Nil},
		{core.Nil, core.Nil},
	} {
		_, ok := core.PApp(EffectSeq, ts...).EvalEffect().(core.ErrorType)
		assert.True(t, ok)
	}
}
