package run

import (
	"github.com/raviqqe/tisp/src/lib/compile"
	"github.com/raviqqe/tisp/src/lib/core"
	"testing"
)

func TestRunWithNoThunk(t *testing.T) {
	Run([]compile.Output{})
}

func TestRunWithOneThunk(t *testing.T) {
	Run([]compile.Output{compile.NewOutput(core.PApp(core.Add, core.NewNumber(123), core.NewNumber(456)), false)})
}

func TestRunWithThunks(t *testing.T) {
	o := compile.NewOutput(core.PApp(core.Add, core.NewNumber(123), core.NewNumber(456)), false)
	Run([]compile.Output{o, o, o, o, o, o, o, o})
}
