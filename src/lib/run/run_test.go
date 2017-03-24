package run

import (
	"testing"

	"github.com/raviqqe/tisp/src/lib/compile"
	"github.com/raviqqe/tisp/src/lib/core"
	"github.com/raviqqe/tisp/src/lib/std"
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
