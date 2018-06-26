package builtins

import (
	"testing"

	"github.com/cloe-lang/cloe/src/lib/core"
	"github.com/stretchr/testify/assert"
)

func TestSeq(t *testing.T) {
	w := core.PApp(Write, core.Nil)

	for _, ts := range [][]core.Value{
		{core.Nil},
		{core.Nil, core.Nil},
		{w},
		{w, w},
	} {
		assert.Equal(
			t,
			core.Nil,
			core.EvalPure(core.PApp(Seq, ts...)))
	}
}

func TestEffectSeq(t *testing.T) {
	assert.Equal(
		t,
		core.Nil,
		core.EvalImpure(core.PApp(
			EffectSeq,
			core.PApp(Write, core.NewNumber(42)),
			core.PApp(Write, core.NewString("OK!")))))
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
