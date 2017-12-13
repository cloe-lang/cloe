package run

import (
	"sync"
	"testing"

	"github.com/coel-lang/coel/src/lib/builtins"
	"github.com/coel-lang/coel/src/lib/compile"
	"github.com/coel-lang/coel/src/lib/core"
	"github.com/stretchr/testify/assert"
)

func TestRunWithNoThunk(t *testing.T) {
	Run(nil)
}

func TestRunWithOneThunk(t *testing.T) {
	Run([]compile.Effect{compile.NewEffect(core.PApp(builtins.Write, core.NewNumber(42)), false)})
}

func TestRunWithThunks(t *testing.T) {
	o := compile.NewEffect(core.PApp(builtins.Write, core.NewNumber(42)), false)
	Run([]compile.Effect{o, o, o, o, o, o, o, o})
}

func TestRunWithExpandedList(t *testing.T) {
	Run([]compile.Effect{compile.NewEffect(
		core.NewList(core.PApp(builtins.Write, core.True), core.PApp(builtins.Write, core.False)),
		true)})
}

func TestEvalEffectListFail(t *testing.T) {
	defer func() {
		assert.NotEqual(t, nil, recover())
	}()

	wg := sync.WaitGroup{}
	wg.Add(1)
	evalEffectList(core.ValueError("Not list effect"), &wg, func(err error) { panic(err) })
}

func TestRunEffectFail(t *testing.T) {
	defer func() {
		assert.NotEqual(t, nil, recover())
	}()

	wg := sync.WaitGroup{}
	sem <- true
	wg.Add(1)
	runEffect(core.Nil, &wg, func(err error) { panic(err) })
}
