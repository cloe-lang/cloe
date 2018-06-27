package builtins

import (
	"testing"

	"github.com/cloe-lang/cloe/src/lib/core"
	"github.com/stretchr/testify/assert"
)

func TestSeq(t *testing.T) {
	for _, ts := range [][]core.Value{
		{core.Nil},
		{core.Nil, core.Nil},
	} {
		assert.Equal(
			t,
			core.Nil,
			core.EvalPure(core.PApp(Seq, ts...)))
	}
}

func TestSeqWithEffects(t *testing.T) {
	w := core.PApp(Print, core.Nil)

	for _, ts := range [][]core.Value{
		{w},
		{w, w},
	} {
		_, ok := core.EvalPure(core.PApp(Seq, ts...)).(*core.ErrorType)
		assert.True(t, ok)
	}
}

func TestEffectSeq(t *testing.T) {
	assert.Equal(
		t,
		core.Nil,
		core.EvalImpure(core.PApp(
			EffectSeq,
			core.PApp(Print, core.NewNumber(42)),
			core.PApp(Print, core.NewString("OK!")))))
}

func TestEffectSeqWithPureValues(t *testing.T) {
	for _, ts := range [][]core.Value{
		{core.Nil},
		{core.Nil, core.Nil},
	} {
		_, ok := core.EvalImpure(core.PApp(EffectSeq, ts...)).(*core.ErrorType)
		assert.True(t, ok)
	}
}
