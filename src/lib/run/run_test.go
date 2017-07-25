package run

import (
	"sync"
	"testing"

	"github.com/tisp-lang/tisp/src/lib/compile"
	"github.com/tisp-lang/tisp/src/lib/core"
	"github.com/tisp-lang/tisp/src/lib/std"
)

func TestRunWithNoThunk(t *testing.T) {
	Run(nil)
}

func TestRunWithOneThunk(t *testing.T) {
	Run([]compile.Output{compile.NewOutput(core.PApp(std.Write, core.NewNumber(42)), false)})
}

func TestRunWithThunks(t *testing.T) {
	o := compile.NewOutput(core.PApp(std.Write, core.NewNumber(42)), false)
	Run([]compile.Output{o, o, o, o, o, o, o, o})
}

func TestRunWithExpandedList(t *testing.T) {
	Run([]compile.Output{compile.NewOutput(
		core.NewList(core.PApp(std.Write, core.True), core.PApp(std.Write, core.False)),
		true)})
}

func TestEvalOutputListFail(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fail()
		}
	}()

	wg := sync.WaitGroup{}
	wg.Add(1)
	evalOutputList(core.ValueError("Not list output"), &wg)
}

func TestRunOutputFail(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fail()
		}
	}()

	wg := sync.WaitGroup{}
	sem <- true
	wg.Add(1)
	runOutput(core.Nil, &wg)
}
